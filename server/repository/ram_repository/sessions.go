package ramrepository

import (
	"httpServer/server/domain"
	"httpServer/server/repository"
)

type SessionsRepository struct {
	data map[string]domain.Session
}

func NewSessionsRepository() (*SessionsRepository, error) {
	return &SessionsRepository{data: make(map[string]domain.Session)}, nil
}

func (rs *SessionsRepository) GetSessionId(sessionId string) (string, error) {
	value, exists := rs.data[sessionId]

	if !exists {
		return "", repository.ErrNotFoundSessionId
	}

	return value.SessionId, nil
}

func (rs *SessionsRepository) PostSessionId(session *domain.Session) (string, error) {
	sessionId := session.SessionId
	rs.data[sessionId] = *session

	return sessionId, nil
}
