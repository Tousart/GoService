package postgres

import (
	"database/sql"
	"errors"
	"httpServer/server/config"
	"httpServer/server/domain"
	"httpServer/server/repository"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(cfg config.Postgres) (*UsersRepository, error) {
	db, err := connectToDB(&cfg.Host, &cfg.Port, &cfg.DBName, &cfg.SSLMode)
	if err != nil {
		return nil, err
	}

	if err := applyMigrations(db); err != nil {
		return nil, err
	}

	return &UsersRepository{
		db: db,
	}, nil
}

func (rs *UsersRepository) PostRegister(user *domain.User) error {
	var exists bool
	err := rs.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", user.Login).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// log.Println("user exists")
		return repository.ErrUserExists
	}

	_, err = rs.db.Exec("INSERT INTO users (login, user_id, password) VALUES ($1, $2, $3)", user.Login, user.Id, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (rs *UsersRepository) PostLogin(login string, password string) (string, error) {
	var (
		id   string
		hash string
	)
	err := rs.db.QueryRow("SELECT user_id, password FROM users WHERE login = $1", login).Scan(&id, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", repository.ErrUserExists
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return "", repository.ErrIncorrectPassword
	}

	return id, nil
}
