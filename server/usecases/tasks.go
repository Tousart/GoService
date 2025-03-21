package usecases

import "server/domain"

type Tasks interface {
	GetStatus(taskId string) (string, error)
	GetResult(taskId string) (string, string, error)
	PostSendTask(translator string, code string) (*domain.Task, error)
	PostCommitTask(taskId string, stdout string, stderr string) error
}
