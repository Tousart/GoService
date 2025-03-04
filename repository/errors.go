package repository

import "errors"

var (
	NotFound          = errors.New("Task id not found")
	NotFoundUser      = errors.New("User not found")
	NotFoundSessionId = errors.New("Session id not found")
	IncorrectPassword = errors.New("Incorrect password")
	UserExists        = errors.New("The user exists")
)
