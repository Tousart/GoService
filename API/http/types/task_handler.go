package types

import (
	"encoding/json"
	"errors"
	"httpServer/repository"
	"httpServer/usecases"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Authorization(r *http.Request, service usecases.Users) error {
	auth_header := r.Header.Get("Authorization")
	if auth_header == "" {
		return errors.New("missing Authorization header")
	} else if !strings.HasPrefix(auth_header, "Bearer ") {
		return errors.New("invalid authorization header format")
	}

	_, err := service.GetSessionId(strings.TrimPrefix(auth_header, "Bearer "))
	if err != nil {
		return errors.New("session id not found")
	}

	return nil
}

// Парсим {task_id}
func CreateGetRequestHandler(r *http.Request) (*GetTaskIdHandler, error) {
	task_id := chi.URLParam(r, "task_id")

	if task_id == "" {
		return nil, errors.New("missing task id")
	}

	return &GetTaskIdHandler{Value: task_id}, nil
}

type GetTaskIdHandler struct {
	Value string `json:"task_id"`
}

type GetResultHandler struct {
	Value string `json:"result"`
}

type GetStatusHandler struct {
	Value string `json:"status"`
}

// Функция вывода ошибки/правильного результата (в зависимости от того, как получили Response)
func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err == repository.NotFound {
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
