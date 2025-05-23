package main

import (
	rabbitmq "httpServer/code_processor/API/rabbitMQ"
	"httpServer/code_processor/config"
	"httpServer/code_processor/metrics"
	"httpServer/code_processor/repository/postgres"
	"httpServer/code_processor/usecases/service"
	"log"
)

func main() {
	codePrcssrFlags := config.ParseFlags()
	var cfg config.CodeProcessorConfig
	config.MustLoad(codePrcssrFlags.CodeProcessorConfigPath, &cfg)

	// Сервер для метрик
	go func() {
		err := metrics.Listen(cfg.Prometheus)
		if err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Осуществляем подключение к бд
	resultRepo, err := postgres.NewResultRepository(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// usecases для сервиса
	codeProcessor, err := service.NewCodeProcessor(cfg.RabbitMQ, resultRepo)
	if err != nil {
		log.Fatalf("Failed to make consumer %v", err)
	}

	// Создаем консьюмер (слушатель очереди)
	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ, codeProcessor)
	if err != nil {
		log.Fatalf("Failed to create consumer %v", err)
	}

	if err = consumer.MakeTask(); err != nil {
		log.Fatalf("Failed to execute task %v", err)
	}

	log.Printf("Starting consumer Code Processor")
}
