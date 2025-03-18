package domain

type Task struct {
	TaskId     string `json:"task_id"`
	Translator string `json:"translator"`
	Code       string `json:"code"`
}
