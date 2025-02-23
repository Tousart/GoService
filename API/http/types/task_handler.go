package types

import (
	"encoding/json"
	"errors"
	"httpServer/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Парсим {task_id}
func CreateGetRequestHandler(r *http.Request) (*GetTaskIdHandler, error) {
	// task_id := r.URL.Query().Get("task_id")

	task_id := chi.URLParam(r, "task_id")

	if task_id == "" {
		return nil, errors.New("Missing Task Id")
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
