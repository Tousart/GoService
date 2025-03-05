package domain

type Task struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
	Result string `json:"result"`
}
