package repository

import "httpServer/domain"

type Tasks interface {
	GetStatus(taskId string) (string, error)
	GetResult(taskId string) (string, error)
	PostTask(task *domain.Task) error
}
