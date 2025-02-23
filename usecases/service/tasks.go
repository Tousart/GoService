package service

import (
	"httpServer/domain"
	"httpServer/repository"
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
func (rs *Tasks) GetStatus(task_id string) (*string, error) {
	return rs.repository.GetStatus(task_id)
}

// Возвращаем результат таски
func (rs *Tasks) GetResult(task_id string) (*string, error) {
	return rs.repository.GetResult(task_id)
}

// Создаем статус и результат таски
func (rs *Tasks) Post(task *domain.Task) error {
	return rs.repository.Post(task)
}
