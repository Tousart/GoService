package rabbitMQ

import (
	"encoding/json"
	"httpServer/domain"

	"github.com/streadway/amqp"
)

type RabbitMQSender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQSender(amqpURL, queueName string) (*RabbitMQSender, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, // имя очереди
		true,      // устойчивость (сохранится при перезапуске сервера)
		false,     // очередь НЕ будет удалена, даже когда нет потребителей
		false,     // будет эксклюзивна только для текущего соединения
		false,     // будет ждать ответа от сервера, что очередь создана
		nil,       // дополнительные аргументы
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQSender{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (r *RabbitMQSender) Send(message *domain.TaskMessage) error {
	body, err := json.Marshal(*message)
	if err != nil {
		return err
	}

	// Отправка сообщения в очередь
	err = r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
