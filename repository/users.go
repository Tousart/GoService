package repository

import "httpServer/domain"

type Users interface {
	PostRegister(*domain.User) error
	PostLogin(login string, password string) (string, error)
}
