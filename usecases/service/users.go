package service

import (
	"crypto/rand"
	"encoding/base64"
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

func (rs *Users) GetSessionId(session_id string) (string, error) {
	return rs.repository.GetSessionId(session_id)
}

func (rs *Users) PostRegister(login string, password string) error {
	user := createUser(login, password)
	return rs.repository.PostRegister(user)
}

func (rs *Users) PostLogin(login string, password string) (string, error) {
	session_id := createToken()
	err := rs.repository.PostLogin(login, password, session_id)
	return session_id, err
}

func createUser(login string, password string) *domain.User {
	return &domain.User{
		Id:       uuid.New().String(),
		Login:    login,
		Password: password,
	}
}

// Гениальная функция генерации токена для сессии
func createToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
