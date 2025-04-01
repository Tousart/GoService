package main

import (
	rabbitmq "httpServer/code_processor/API/rabbitMQ"
	"httpServer/code_processor/config"
	"httpServer/code_processor/metrics"
	"httpServer/code_processor/repository/postgres"
	"httpServer/code_processor/service"
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

	// "amqp://guest:guest@rabbitMQ:5672"
	codeProcessor, err := service.NewCodeProcessor(cfg.RabbitMQ, resultRepo)
	if err != nil {
		log.Fatalf("Failed to make consumer %v", err)
	}

	consumer := rabbitmq.NewConsumer(*codeProcessor)
	if err = consumer.MakeTask(); err != nil {
		log.Fatalf("Failed to execute task %v", err)
	}

	log.Printf("Starting consumer Code Processor")
}
