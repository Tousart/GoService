package usecases

type Sessions interface {
	GetSessionId(session_id string) (string, error)
	PostSessionId(user_id string) (string, error)
}
