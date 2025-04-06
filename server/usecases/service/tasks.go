package service

import (
	"httpServer/server/domain"
	"httpServer/server/repository"

	"github.com/google/uuid"
)

type Tasks struct {
	repository repository.Tasks
	sender     repository.TaksSender
}

func NewTasks(repo repository.Tasks,
	sender repository.TaksSender) *Tasks {
	return &Tasks{
		repository: repo,
		sender:     sender,
	}
}

func (rs *Tasks) GetStatus(taskId string) (string, error) {
	return rs.repository.GetStatus(taskId)
}

func (rs *Tasks) GetResult(taskId string) (string, string, error) {
	return rs.repository.GetResult(taskId)
}

func (rs *Tasks) PostSendTask(translator string, code string) (*domain.Task, error) {
	task := createTask(uuid.New().String(), "in_progress", "", "")
	rs.repository.PostTask(task)

	taskMessage := createTaskMessage(task.TaskId, translator, code)
	err := rs.sender.Send(taskMessage)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func createTask(id string, status string, stdout string, stderr string) *domain.Task {
	return &domain.Task{
		TaskId: id,
		Status: status,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func createTaskMessage(id string, trans string, code string) *domain.TaskMessage {
	return &domain.TaskMessage{
		TaskId:     id,
		Translator: trans,
		Code:       code,
	}
}
