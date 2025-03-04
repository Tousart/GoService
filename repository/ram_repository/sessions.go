package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

type SessionsRepository struct {
	data map[string]domain.Session
}

func NewSessionsRepository() *SessionsRepository {
	return &SessionsRepository{data: make(map[string]domain.Session)}
}

func (rs *SessionsRepository) GetSessionId(session_id string) (string, error) {
	value, exists := rs.data[session_id]

	if !exists {
		return "", repository.ErrNotFoundSessionId
	}

	return value.Session_id, nil
}

func (rs *SessionsRepository) PostSessionId(session *domain.Session) (string, error) {
	session_id := session.Session_id
	user_id := session.User_id

	rs.data[session.Session_id] = domain.Session{
		Session_id: session_id,
		User_id:    user_id,
	}

	return session_id, nil
}
