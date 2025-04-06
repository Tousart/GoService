package types

import (
	"encoding/json"
	"errors"
	"httpServer/server/usecases"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Authorization(r *http.Request, service usecases.Sessions) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("missing Authorization header")
	} else if !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("invalid authorization header format")
	}

	_, err := service.GetSessionId(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		return errors.New("session id not found")
	}

	return nil
}

func CreateGetRequestHandler(r *http.Request) (*GetTaskIdHandler, error) {
	taskId := chi.URLParam(r, "task_id")

	if taskId == "" {
		return nil, errors.New("missing task id")
	}

	return &GetTaskIdHandler{Value: taskId}, nil
}

func CreateTaskRequestHandler(r *http.Request) (*TaskBody, error) {
	var taskBody TaskBody
	err := json.NewDecoder(r.Body).Decode(&taskBody)

	if err != nil {
		return nil, errors.New("bad request")
	}

	return &taskBody, nil
}

type GetTaskIdHandler struct {
	Value string `json:"task_id"`
}

type GetResultHandler struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type GetStatusHandler struct {
	Value string `json:"status"`
}

type TaskBody struct {
	Translator string `json:"translator"`
	Code       string `json:"code"`
}
