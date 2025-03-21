package ramrepository

import (
	"httpServer/domain"
	"httpServer/repository"
)

type Tasks struct {
	data map[string]domain.Task
}

func NewTasks() *Tasks {
	return &Tasks{
		data: make(map[string]domain.Task),
	}
}

// Возвращаем статус таски
func (rs *Tasks) GetStatus(taskId string) (string, error) {
	value, exists := rs.data[taskId]

	if !exists {
		return "", repository.ErrNotFound
	}

	return value.Status, nil
}

// Возвращаем результат таски
func (rs *Tasks) GetResult(taskId string) (string, string, error) {
	value, exists := rs.data[taskId]

	if !exists {
		return "", "", repository.ErrNotFound
	}

	return value.Stdout, value.Stderr, nil
}

// Создаем статус и результат таски
func (rs *Tasks) PostTask(task *domain.Task) error {
	rs.data[task.TaskId] = domain.Task{
		Status: task.Status,
		Stdout: task.Stdout,
		Stderr: task.Stderr,
	}

	return nil
}
