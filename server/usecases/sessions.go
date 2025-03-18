package usecases

type Sessions interface {
	GetSessionId(sessionId string) (string, error)
	PostSessionId(userId string) (string, error)
}
