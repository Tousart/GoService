package domain

type TaskMessage struct {
	TaskId     string `json:"task_id"`
	Translator string `json:"translator"`
	Code       string `json:"code"`
}
