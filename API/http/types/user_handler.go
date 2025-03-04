package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"httpServer/repository"
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

func ProcessErrorRegister(w http.ResponseWriter, err error, resp string) {
	if err == repository.UserExists {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Пользователь %s зарегистрирован.", resp)
}

func ProcessErrorLogin(w http.ResponseWriter, err error, resp any) {
	if err == repository.NotFound {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err == repository.IncorrectPassword {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
