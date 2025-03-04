package usecases

type Users interface {
	GetSessionId(session_id string) (string, error)
	PostRegister(login string, password string) error
	PostLogin(login string, password string) (string, error)
}
