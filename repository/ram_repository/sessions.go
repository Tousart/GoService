package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

type SessionsRepository struct {
	data      map[string]domain.Session
	userToken map[string]string
}

func NewSessionsRepository() *SessionsRepository {
	return &SessionsRepository{
		data:      make(map[string]domain.Session),
		userToken: make(map[string]string),
	}
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
	userId := session.UserId

	if id, exists := rs.userToken[userId]; exists {
		delete(rs.data, id)
	}

	rs.data[sessionId] = *session
	rs.userToken[userId] = sessionId

	return sessionId, nil
}
