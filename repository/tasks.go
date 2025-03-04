package repository

import "httpServer/domain"

type Tasks interface {
	GetStatus(task_id string) (string, error)
	GetResult(task_id string) (string, error)
	PostTask(task *domain.Task) error
}
