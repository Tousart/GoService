package types

import (
	"encoding/json"
	"errors"
	"net/http"
)

func CreateUserRequestHandler(r *http.Request) (*User, error) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		return &user, errors.New("bad request")
	}

	return &user, nil
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetSessionIdHandler struct {
	SessionId string `json:"token"`
}
