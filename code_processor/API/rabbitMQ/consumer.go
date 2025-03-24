package rabbitmq

import (
	"encoding/json"
	"httpServer/code_processor/domain"
	"httpServer/code_processor/service"
	"log"
)

type Consumer struct {
	Consumer service.CodeProcessor
}

func NewConsumer(codeProcessor service.CodeProcessor) *Consumer {
	return &Consumer{Consumer: codeProcessor}
}

func (cp *Consumer) MakeTask() error {
	msgs, err := cp.Consumer.Channel.Consume(
		cp.Consumer.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// Берем значения полей для таски из тела сообщения
			var task domain.Task
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Printf("failed to unmarshal body: %v", err)
				continue
			}
			// log.Printf("%s, %s, %s\n", task.TaskId, task.Translator, task.Code)

			// Выполняем код
			stdout, stderr, err := cp.Consumer.ExecuteCodeInDocker(task)
			if err != nil {
				log.Printf("Failed to execute code: %v", err)
				continue
			}

			err = cp.Consumer.SendResult(task.TaskId, stdout, stderr)
			if err != nil {
				log.Printf("Failed to send result: %v", err)
				continue
			}
		}
	}()
	<-forever

	return nil
}
