package service

import (
	"bytes"
	"context"
	"fmt"
	"httpServer/code_processor/config"
	"httpServer/code_processor/domain"
	"httpServer/code_processor/repository"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	_ "github.com/lib/pq"
)

type CodeProcessor struct {
	dockerCli *client.Client
	db        repository.Result
}

func NewCodeProcessor(cfg config.RabbitMQ, db repository.Result) (*CodeProcessor, error) {
	// Создаем Docker-клиент
	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &CodeProcessor{
		dockerCli: dockerCli,
		db:        db,
	}, nil
}

func (cp *CodeProcessor) ExecuteCodeInDocker(task domain.Task) (stdout, stderr string, err error) {
	// Записываем метрику времени выполнения запроса
	startTime := time.Now()
	defer func() {
		observeReqDuration(time.Since(startTime), task.Translator)
		observeUsedTranslator(task.Translator)
	}()

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
	reader, err := cp.dockerCli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to pull image: %v", err)
	}
	defer reader.Close()

	// Ждем завершения загрузки образа
	io.Copy(io.Discard, reader)

	/* Cоздание контейнера и получение результата выполнения кода */

	// Создаем контейнер
	resp, err := cp.dockerCli.ContainerCreate(ctx,
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
	err = cp.dockerCli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	// Таймаут для выполнения задачи
	timeout := 30 * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Ждем завершения контейнера
	statusCh, errCh := cp.dockerCli.ContainerWait(timeoutCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", "", fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Получаем логи (stdout и stderr)
	out, err := cp.dockerCli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
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
	err = cp.dockerCli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
		Force: true, // Принудительное удаление, если контейнер еще работает
	})
	if err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}

func (cp *CodeProcessor) SendResult(taskId string, stdout string, stderr string) error {
	result := domain.Result{
		TaskId: taskId,
		Status: "ready",
		Stdout: stdout,
		Stderr: stderr,
	}

	// Отправляем результат в бд
	cp.db.SendResult(&result)

	return nil
}
