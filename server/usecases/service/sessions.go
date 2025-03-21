package service

import (
	"crypto/rand"
	"encoding/base64"
	"server/domain"
	"server/repository"
)

type Sessions struct {
	repository repository.Sessions
}

func NewSessions(repo repository.Sessions) *Sessions {
	return &Sessions{
		repository: repo,
	}
}

func (rs *Sessions) GetSessionId(sessionId string) (string, error) {
	return rs.repository.GetSessionId(sessionId)
}

func (rs *Sessions) PostSessionId(userId string) (string, error) {
	sessionId := createToken()
	return rs.repository.PostSessionId(&domain.Session{SessionId: sessionId, UserId: userId})
}

// Гениальная функция генерации токена для сессии
func createToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
