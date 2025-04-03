package rabbitmq

import (
	"encoding/json"
	"fmt"
	"httpServer/code_processor/config"
	"httpServer/code_processor/domain"
	"httpServer/code_processor/usecases"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	consumer   usecases.CodeProcessor
}

func NewConsumer(cfg config.RabbitMQ, codeProcessor usecases.CodeProcessor) (*Consumer, error) {
	amqpURL := fmt.Sprintf("amqp://guest:guest@%s:%d", cfg.Host, cfg.Port)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.QueueName, // имя очереди
		true,          // устойчивость (сохранится при перезапуске сервера)
		false,         // очередь НЕ будет удалена, даже когда нет потребителей
		false,         // будет эксклюзивна только для текущего соединения
		false,         // будет ждать ответа от сервера, что очередь создана
		nil,           // дополнительные аргументы
	)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer:   codeProcessor,
		connection: conn,
		channel:    ch,
		queueName:  cfg.QueueName,
	}, nil
}

func (cp *Consumer) MakeTask() error {
	msgs, err := cp.channel.Consume(
		cp.queueName,
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
			stdout, stderr, err := cp.consumer.ExecuteCodeInDocker(task)
			if err != nil {
				log.Printf("Failed to execute code: %v", err)
				continue
			}

			err = cp.consumer.SendResult(task.TaskId, stdout, stderr)
			if err != nil {
				log.Printf("Failed to send result: %v", err)
				continue
			}
		}
	}()
	<-forever

	return nil
}
