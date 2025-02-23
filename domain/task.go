package domain

import "github.com/google/uuid"

type Task struct {
	Task_id string `json:"task_id"`
	Status  string `json:"status"`
	Result  string `json:"result"`
}

func CreateTask() *Task {
	return &Task{
		Task_id: uuid.New().String(),
		Status:  "in_progress",
		Result:  "nothing",
	}
}
