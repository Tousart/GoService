package repository

import "errors"

var (
	ErrNotFound          = errors.New("task id not found")
	ErrNotFoundUser      = errors.New("user not found")
	ErrNotFoundSessionId = errors.New("session id not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserExists        = errors.New("the user exists")
)
