package repository

import "httpServer/server/domain"

type Sessions interface {
	GetSessionId(sessionIdd string) (string, error)
	PostSessionId(session *domain.Session) (string, error)
}
