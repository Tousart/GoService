package domain

type User struct {
	Id       string `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
