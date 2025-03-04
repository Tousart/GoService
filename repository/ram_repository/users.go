package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

type UsersRepository struct {
	data map[string]domain.User
}

type SessionsRepository struct {
	data map[string]domain.Session
}

type Users struct {
	users    *UsersRepository
	sessions *SessionsRepository
}

func NewUsersRepository() *Users {
	return &Users{
		users:    &UsersRepository{data: make(map[string]domain.User)},
		sessions: &SessionsRepository{data: make(map[string]domain.Session)},
	}
}

func (rs *Users) GetSessionId(session_id string) (string, error) {
	value, exists := rs.sessions.data[session_id]

	if !exists {
		return "", repository.NotFoundSessionId
	}

	return value.Session_id, nil
}

func (rs *Users) PostRegister(user *domain.User) error {
	_, exists := rs.users.data[user.Login]

	if exists {
		return repository.UserExists
	}

	rs.users.data[user.Login] = *user

	return nil
}

func (rs *Users) PostLogin(login string, password string, session_id string) error {
	value, exists := rs.users.data[login]

	if !exists {
		return repository.NotFoundUser
	} else if value.Password != password {
		return repository.IncorrectPassword
	}

	rs.sessions.data[session_id] = domain.Session{
		Session_id: session_id,
		User_id:    value.Id,
	}

	return nil
}
