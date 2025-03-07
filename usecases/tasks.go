package usecases

import "httpServer/domain"

type Tasks interface {
	GetStatus(taskId string) (string, error)
	GetResult(taskId string) (string, error)
	PostTask() (*domain.Task, error)
}
