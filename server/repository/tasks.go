package repository

import "server/domain"

type Tasks interface {
	GetStatus(taskId string) (string, error)
	GetResult(taskId string) (string, string, error)
	PostTask(task *domain.Task) error
}
