package service

import (
	"httpServer/domain"
	"httpServer/repository"
	"time"

	"github.com/google/uuid"
)

type Tasks struct {
	repository repository.Tasks
}

func NewTasks(repo repository.Tasks) *Tasks {
	return &Tasks{
		repository: repo,
	}
}

// Возвращаем статус таски
func (rs *Tasks) GetStatus(task_id string) (string, error) {
	return rs.repository.GetStatus(task_id)
}

// Возвращаем результат таски
func (rs *Tasks) GetResult(task_id string) (string, error) {
	return rs.repository.GetResult(task_id)
}

// Создаем статус и результат таски
func (rs *Tasks) PostTask() (*domain.Task, error) {
	task := createTask()
	rs.repository.PostTask(task)

	// имитация бурной деятельности
	go func(task *domain.Task) {
		time.Sleep(40 * time.Second)

		task.Status = "ready"
		task.Result = "end"

		rs.repository.PostTask(task)
	}(task)

	return task, nil
}

func createTask() *domain.Task {
	return &domain.Task{
		Task_id: uuid.New().String(),
		Status:  "in_progress",
		Result:  "nothing",
	}
}
