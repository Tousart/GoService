package usecases

import "httpServer/server/domain"

type Tasks interface {
	GetStatus(taskId string) (string, error)
	GetResult(taskId string) (string, string, error)
	PostSendTask(translator string, code string) (*domain.Task, error)
	// PostCommitTask(taskId string, status string, stdout string, stderr string) error
}
