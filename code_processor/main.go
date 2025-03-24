package main

import (
	rabbitmq "httpServer/code_processor/API/rabbitMQ"
	"httpServer/code_processor/config"
	"httpServer/code_processor/service"
	"log"
)

func main() {
	codePrcssrFlags := config.ParseFlags()
	var cfg config.CodeProcessorConfig
	config.MustLoad(codePrcssrFlags.CodeProcessorConfigPath, &cfg)

	// "amqp://guest:guest@rabbitMQ:5672"
	codeProcessor, err := service.NewCodeProcessor(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to make consumer %v", err)
	}

	consumer := rabbitmq.NewConsumer(*codeProcessor)
	if err = consumer.MakeTask(); err != nil {
		log.Fatalf("Failed to execute task %v", err)
	}
}
