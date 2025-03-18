package main

import (
	"code_processor/service"
	"log"
)

func main() {
	codeProcessor, err := service.NewCodeProcessor("amqp://guest:guest@rabbitMQ:5672", "testQueue")
	if err != nil {
		log.Fatalf("Failed to make code processor %v", err)
	}

	err = codeProcessor.MakeTask()
	if err != nil {
		log.Fatalf("Error processing tasks: %v", err)
	}
}
