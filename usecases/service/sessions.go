package service

import (
	"crypto/rand"
	"encoding/base64"
	"httpServer/domain"
	"httpServer/repository"
)

type Sessions struct {
	repository repository.Sessions
}

func NewSessions(repo repository.Sessions) *Sessions {
	return &Sessions{
		repository: repo,
	}
}

func (rs *Sessions) GetSessionId(session_id string) (string, error) {
	return rs.repository.GetSessionId(session_id)
}

func (rs *Sessions) PostSessionId(user_id string) (string, error) {
	session_id := createToken()
	return rs.repository.PostSessionId(&domain.Session{Session_id: session_id, User_id: user_id})
}

// Гениальная функция генерации токена для сессии
func createToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
