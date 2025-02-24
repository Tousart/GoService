package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

// Хранилище - мапа с тасками по айдишникам
type Tasks struct {
	data map[string]domain.Task
}

func NewTasks() *Tasks {
	return &Tasks{
		data: make(map[string]domain.Task),
	}
}

// Возвращаем статус таски
func (rs *Tasks) GetStatus(task_id string) (string, error) {
	value, exists := rs.data[task_id]

	if !exists {
		return "", repository.NotFound
	}

	return value.Status, nil
}

// Возвращаем результат таски
func (rs *Tasks) GetResult(task_id string) (string, error) {
	value, exists := rs.data[task_id]

	if !exists {
		return "", repository.NotFound
	}

	return value.Result, nil
}

// Создаем статус и результат таски
func (rs *Tasks) Post(task *domain.Task) error {
	rs.data[task.Task_id] = domain.Task{
		Status: task.Status,
		Result: task.Result,
	}

	return nil
}
