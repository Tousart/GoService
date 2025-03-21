package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	rabbitmq "httpServer/API/rabbitMQ"
	"httpServer/domain"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stdcopy"
)

type CodeProcessor struct {
	Consumer *rabbitmq.Consumer
}

func NewCodeProcessor(consumer *rabbitmq.Consumer) *CodeProcessor {
	return &CodeProcessor{Consumer: consumer}
}

func (cp *CodeProcessor) MakeTask() error {
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
			stdout, stderr, err := cp.executeCodeInDocker(task)
			if err != nil {
				log.Printf("Failed to execute code: %v", err)
				continue
			}

			// Кладем результат кода в json для отправки ответа
			resultBody, err := json.Marshal(domain.Result{
				TaskId: task.TaskId,
				Stdout: stdout,
				Stderr: stderr,
			})
			if err != nil {
				log.Printf("Failed to marshal result: %v", err)
				continue
			}
			// log.Printf("Result %s", string(resultBody))

			// Создаем и делаем запрос на отправку данных в бд
			req, err := http.NewRequest("POST", "http://http_server:8080/commit", bytes.NewBuffer([]byte(resultBody)))
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				continue
			}

			client := &http.Client{}
			_, err = client.Do(req)
			if err != nil {
				log.Printf("Failed to make request: %v", err)
				continue
			}
			// log.Println("Request completed")
		}
	}()
	<-forever

	return nil
}

func (cp *CodeProcessor) executeCodeInDocker(task domain.Task) (stdout, stderr string, err error) {
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
	reader, err := cp.Consumer.DockerCli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to pull image: %v", err)
	}
	defer reader.Close()

	// Ждем завершения загрузки образа
	io.Copy(io.Discard, reader)

	/* Cоздание контейнера и получение результата выполнения кода */

	// Создаем контейнер
	resp, err := cp.Consumer.DockerCli.ContainerCreate(ctx,
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
	err = cp.Consumer.DockerCli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	// Таймаут для выполнения задачи
	timeout := 30 * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Ждем завершения контейнера
	statusCh, errCh := cp.Consumer.DockerCli.ContainerWait(timeoutCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", "", fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Получаем логи (stdout и stderr)
	out, err := cp.Consumer.DockerCli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
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
	err = cp.Consumer.DockerCli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
		Force: true, // Принудительное удаление, если контейнер еще работает
	})
	if err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
