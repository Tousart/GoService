package repository

import "httpServer/server/domain"

type Users interface {
	PostRegister(*domain.User) error
	PostLogin(login string, password string) (string, error)
}
