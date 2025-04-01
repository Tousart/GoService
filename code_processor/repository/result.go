package repository

import "httpServer/code_processor/domain"

type Result interface {
	SendResult(result *domain.Result) error
}
