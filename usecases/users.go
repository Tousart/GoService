package usecases

type Users interface {
	PostRegister(login string, password string) error
	PostLogin(login string, password string) (string, error)
}
