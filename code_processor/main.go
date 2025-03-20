package main

import (
	rabbitmq "code_processor/API/rabbitMQ"
	"code_processor/config"
	"code_processor/service"
	"log"
)

func main() {
	codePrcssrFlags := config.ParseFlags()
	var cfg config.CodeProcessorConfig
	config.MustLoad(codePrcssrFlags.CodeProcessorConfigPath, &cfg)

	// "amqp://guest:guest@rabbitMQ:5672"
	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to make consumer %v", err)
	}

	codeProcessor := service.NewCodeProcessor(consumer)
	if err = codeProcessor.MakeTask(); err != nil {
		log.Fatalf("Failed to execute task %v", err)
	}
}
