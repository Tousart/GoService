package rabbitmq

import (
	"code_processor/config"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
	DockerCli  *client.Client
}

func NewConsumer(cfg config.RabbitMQ) (*Consumer, error) {
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

	// Создаем Docker-клиент
	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &Consumer{
		Connection: conn,
		Channel:    ch,
		QueueName:  cfg.QueueName,
		DockerCli:  dockerCli,
	}, nil
}
