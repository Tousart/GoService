package repository

import "server/domain"

type TaksSender interface {
	Send(message *domain.TaskMessage) error
}
