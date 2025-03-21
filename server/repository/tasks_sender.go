package repository

import "httpServer/server/domain"

type TaksSender interface {
	Send(message *domain.TaskMessage) error
}
