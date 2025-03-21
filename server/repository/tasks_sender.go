package repository

import "httpServer/domain"

type TaksSender interface {
	Send(message *domain.TaskMessage) error
}
