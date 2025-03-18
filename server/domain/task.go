package domain

type Task struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
