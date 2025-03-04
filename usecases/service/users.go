package service

import (
	"httpServer/domain"
	"httpServer/repository"

	"github.com/google/uuid"
)

type Users struct {
	repository repository.Users
}

func NewUsers(repo repository.Users) *Users {
	return &Users{
		repository: repo,
	}
}

func (rs *Users) PostRegister(login string, password string) error {
	user := createUser(login, password)
	return rs.repository.PostRegister(user)
}

func (rs *Users) PostLogin(login string, password string) (string, error) {
	return rs.repository.PostLogin(login, password)
}

func createUser(login string, password string) *domain.User {
	return &domain.User{
		Id:       uuid.New().String(),
		Login:    login,
		Password: password,
	}
}
