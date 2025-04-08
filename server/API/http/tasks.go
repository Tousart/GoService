package http

import (
	"net/http"

	"httpServer/server/API/http/types"
	"httpServer/server/usecases"

	"github.com/go-chi/chi/v5"
)

type Tasks struct {
	serviceTasks    usecases.Tasks
	serviceSessions usecases.Sessions
}

func NewTasksHandler(tasks usecases.Tasks, sessions usecases.Sessions) *Tasks {
	return &Tasks{
		serviceTasks:    tasks,
		serviceSessions: sessions,
	}
}

// @Summary Get a Status
// @Description Get a tasks status by its id
// @Tags status
// @Accept json
// @Produce json
// @Param task_id path string true "Task Id"
// @Param Authorization header string true "Bearer {auth_token}"
// @Success 200 {object} types.GetStatusHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorize"
// @Failure 404 {string} string "Task id not found"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /status/{task_id} [get]
func (s *Tasks) getHandlerStatus(w http.ResponseWriter, r *http.Request) {
	if err := types.Authorization(r, s.serviceSessions); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := types.CreateGetRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	status, err := s.serviceTasks.GetStatus(req.Value)
	types.ProcessErrorGetTask(w, err, &types.GetStatusHandler{Value: status})
}

// @Summary Get a Result
// @Description Get a tasks result by its id
// @Tags result
// @Accept json
// @Produce json
// @Param task_id path string true "Task Id"
// @Param Authorization header string true "Bearer {auth_token}"
// @Success 200 {object} types.GetResultHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorize"
// @Failure 404 {string} string "Task id not found"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /result/{task_id} [get]
func (s *Tasks) getHandlerResult(w http.ResponseWriter, r *http.Request) {
	if err := types.Authorization(r, s.serviceSessions); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := types.CreateGetRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	stdout, stderr, err := s.serviceTasks.GetResult(req.Value)
	types.ProcessErrorGetTask(w, err, &types.GetResultHandler{Stdout: stdout, Stderr: stderr})
}

// @Summary Post a Task
// @Description make task and get tsk id
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {auth_token}"
// @Param taskBody body types.TaskBody true "TaskBody"
// @Success 201 {object} types.GetTaskIdHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorize"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /task [post]
func (s *Tasks) postHandlerCreate(w http.ResponseWriter, r *http.Request) {
	if err := types.Authorization(r, s.serviceSessions); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := types.CreateTaskRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task, err := s.serviceTasks.PostSendTask(req.Translator, req.Code)
	types.ProcessErrorCreateTask(w, err, &types.GetTaskIdHandler{Value: task.TaskId})
}

// func (s *Tasks) postHandlerCommit(w http.ResponseWriter, r *http.Request) {
// 	req, err := types.CreateTaskCommitHandler(r)
// 	if err != nil {
// 		http.Error(w, "Bad Request", http.StatusBadRequest)
// 		return
// 	}

// 	err = s.serviceTasks.PostCommitTask(req.TaskId, req.Status, req.Stdout, req.Stderr)
// 	if err != nil {
// 		http.Error(w, "Internal Error", http.StatusInternalServerError)
// 		return
// 	}
// }

func (s *Tasks) WithTasksHandlers(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/status/{task_id}", s.getHandlerStatus)
		r.Get("/result/{task_id}", s.getHandlerResult)
		r.Post("/task", s.postHandlerCreate)
	})
}
