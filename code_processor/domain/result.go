package domain

type Result struct {
	TaskId string `json:"task_id"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
