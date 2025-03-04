package repository

import "httpServer/domain"

type Users interface {
	GetSessionId(sessionId string) (string, error)
	PostRegister(*domain.User) error
	PostLogin(login string, password string, sessionId string) error
}
