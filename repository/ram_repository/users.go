package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

type UsersRepository struct {
	data map[string]domain.User
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{data: make(map[string]domain.User)}
}

func (rs *UsersRepository) PostRegister(user *domain.User) error {
	_, exists := rs.data[user.Login]

	if exists {
		return repository.ErrUserExists
	}

	rs.data[user.Login] = *user

	return nil
}

func (rs *UsersRepository) PostLogin(login string, password string) (string, error) {
	value, exists := rs.data[login]

	if !exists {
		return "", repository.ErrNotFoundUser
	} else if value.Password != password {
		return "", repository.ErrIncorrectPassword
	}

	return value.Id, nil
}
