package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"httpServer/code_processor/config"
	"httpServer/code_processor/domain"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/streadway/amqp"
)

type CodeProcessor struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
	DockerCli  *client.Client
}

func NewCodeProcessor(cfg config.RabbitMQ) (*CodeProcessor, error) {
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

	return &CodeProcessor{
		Connection: conn,
		Channel:    ch,
		QueueName:  cfg.QueueName,
		DockerCli:  dockerCli,
	}, nil
}

func (cp *CodeProcessor) ExecuteCodeInDocker(task domain.Task) (stdout, stderr string, err error) {
	// Определяем образ
	var imageName string
	if task.Translator == "python3" {
		imageName = "python:3.9-alpine"
	} else if task.Translator == "gcc" {
		imageName = "gcc:latest"
	} else if task.Translator == "clang" {
		imageName = "clang:latest"
	} else {
		return "", "", fmt.Errorf("unsupported translator: %s", task.Translator)
	}

	// Создаем команду, которая выполнит код пользователя
	var cmd []string
	if task.Translator == "python3" {
		cmd = []string{"python3", "-c", task.Code}
	} else if task.Translator == "gcc" || task.Translator == "clang" {
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.c && %s /tmp/code.c -o /tmp/out && /tmp/out", task.Code, task.Translator)}
	}

	ctx := context.Background()

	// Загружаем образ
	reader, err := cp.DockerCli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to pull image: %v", err)
	}
	defer reader.Close()

	// Ждем завершения загрузки образа
	io.Copy(io.Discard, reader)

	/* Cоздание контейнера и получение результата выполнения кода */

	// Создаем контейнер
	resp, err := cp.DockerCli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			Cmd:   cmd,
			Tty:   false,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:   100 * 1024 * 1024, // Ограничение памяти: 100 MB
				CPUQuota: 50000,             // Ограничение CPU: 50% (в единицах cgroups)
			},
			AutoRemove: false,
		},
		nil,
		nil,
		"")
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %v", err)
	}

	// Запускаем контейнер
	err = cp.DockerCli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	// Таймаут для выполнения задачи
	timeout := 30 * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Ждем завершения контейнера
	statusCh, errCh := cp.DockerCli.ContainerWait(timeoutCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", "", fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Получаем логи (stdout и stderr)
	out, err := cp.DockerCli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer out.Close()

	// Читаем логи
	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, out)
	if err != nil {
		return "", "", fmt.Errorf("failed to read container logs: %v", err)
	}

	// Удаляем контейнер вручную
	err = cp.DockerCli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
		Force: true, // Принудительное удаление, если контейнер еще работает
	})
	if err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}

func (cp *CodeProcessor) SendResult(taskId string, stdout string, stderr string) error {
	// Кладем результат кода в json для отправки ответа
	resultBody, err := json.Marshal(domain.Result{
		TaskId: taskId,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		log.Printf("Failed to marshal result: %v", err)
		return err
	}
	// log.Printf("Result %s", string(resultBody))

	// Создаем и делаем запрос на отправку данных в бд
	req, err := http.NewRequest("POST", "http://http_server:8080/commit", bytes.NewBuffer([]byte(resultBody)))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return err
	}
	// log.Println("Request completed")
	return nil
}
