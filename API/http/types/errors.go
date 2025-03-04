package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"httpServer/repository"
	"net/http"
)

func ProcessErrorTask(w http.ResponseWriter, err error, resp any) {
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "Task id not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}

func ProcessErrorRegister(w http.ResponseWriter, err error, resp string) {
	if errors.Is(err, repository.ErrUserExists) {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Пользователь %s зарегистрирован.", resp)
}

func ProcessErrorLogin(w http.ResponseWriter, err error, resp any) {
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if errors.Is(err, repository.ErrIncorrectPassword) {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
