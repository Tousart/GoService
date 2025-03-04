package repository

import "httpServer/domain"

type Sessions interface {
	GetSessionId(session_id string) (string, error)
	PostSessionId(session *domain.Session) (string, error)
}
