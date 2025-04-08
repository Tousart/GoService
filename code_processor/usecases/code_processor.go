package usecases

import "httpServer/code_processor/domain"

type CodeProcessor interface {
	ExecuteCodeInDocker(task domain.Task) (stdout, stderr string, err error)
	SendResult(taskId string, stdout string, stderr string) error
}
