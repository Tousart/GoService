package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

// Таска состоит из статуса и результата
type Tasks_values struct {
	status string
	result string
}

// Хранилище - мапа с тасками по айдишникам
type Tasks struct {
	data map[string]Tasks_values
}

func NewTasks() *Tasks {
	return &Tasks{
		data: make(map[string]Tasks_values),
	}
}

// Возвращаем статус таски
func (rs *Tasks) GetStatus(task_id string) (*string, error) {
	value, exists := rs.data[task_id]

	if !exists {
		return nil, repository.NotFound
	}

	return &value.status, nil
}

// Возвращаем результат таски
func (rs *Tasks) GetResult(task_id string) (*string, error) {
	value, exists := rs.data[task_id]

	if !exists {
		return nil, repository.NotFound
	}

	return &value.result, nil
}

// Создаем статус и результат таски
func (rs *Tasks) Post(task *domain.Task) error {
	rs.data[task.Task_id] = Tasks_values{
		status: task.Status,
		result: task.Result,
	}

	return nil
}
