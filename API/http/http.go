package http

import (
	"net/http"
	"time"

	"httpServer/API/http/types"
	"httpServer/domain"
	"httpServer/usecases"

	"github.com/go-chi/chi/v5"
)

type Tasks struct {
	service usecases.Tasks
}

func NewTasksHandler(service usecases.Tasks) *Tasks {
	return &Tasks{service: service}
}

// @Summary Get a Status
// @Description Get a tasks status by its id
// @Tags status
// @Accept json
// @Produce json
// @Param task_id path string true "Task Id"
// @Success 200 {object} types.GetStatusHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /status/{task_id} [get]
func (s *Tasks) getHandlerStatus(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	value, err := s.service.GetStatus(req.Value)
	types.ProcessError(w, err, &types.GetStatusHandler{Value: *value})
}

// @Summary Get a Result
// @Description Get a tasks result by its id
// @Tags result
// @Accept json
// @Produce json
// @Param task_id path string true "Task Id"
// @Success 200 {object} types.GetResultHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /result/{task_id} [get]
func (s *Tasks) getHandlerResult(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	value, err := s.service.GetResult(req.Value)
	types.ProcessError(w, err, &types.GetResultHandler{Value: *value})
}

// @Summary Post a Task
// @Description make task and get tsk id
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} types.GetTaskIdHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [post]
func (s *Tasks) postHandler(w http.ResponseWriter, r *http.Request) {
	task := domain.CreateTask()
	err := s.service.Post(task)
	types.ProcessError(w, err, &types.GetTaskIdHandler{Value: task.Task_id})

	// имитация бурной деятельности
	go func(task *domain.Task) {
		time.Sleep(40 * time.Second)

		task.Status = "ready"
		task.Result = "end"
		s.service.Post(task)

	}(task)
}

func (s *Tasks) WithTasksHandlers(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/status/{task_id}", s.getHandlerStatus)
		r.Get("/result/{task_id}", s.getHandlerResult)
		r.Post("/task", s.postHandler)
	})
}
