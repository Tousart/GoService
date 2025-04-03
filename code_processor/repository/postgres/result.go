package postgres

import (
	"database/sql"
	"httpServer/code_processor/config"
	"httpServer/code_processor/domain"
)

type ResultRepository struct {
	db *sql.DB
}

func NewResultRepository(cfg config.Postgres) (*ResultRepository, error) {
	db, err := connectToDB(&cfg.Host, &cfg.Port, &cfg.DBName, &cfg.SSLMode)
	if err != nil {
		return nil, err
	}

	return &ResultRepository{
		db: db,
	}, nil
}

func (rs *ResultRepository) SendResult(result *domain.Result) error {
	_, err := rs.db.Exec("UPDATE tasks SET status = $1, stdout = $2, stderr = $3 WHERE task_id = $4", result.Status, result.Stdout, result.Stderr, result.TaskId)
	if err != nil {
		return err
	}

	// Кладем результат кода в json для отправки ответа
	// resultJson, err := json.Marshal(result)
	// if err != nil {
	// 	log.Printf("Failed to marshal result: %v", err)
	// 	return err
	// }
	// log.Printf("Result %s", string(resultJson))

	// Создаем и делаем запрос на отправку данных в бд
	// req, err := http.NewRequest("POST", "http://http_server:8080/commit", bytes.NewBuffer([]byte(resultJson)))
	// if err != nil {
	// 	log.Printf("Failed to create request: %v", err)
	// 	return err
	// }

	// client := &http.Client{}
	// _, err = client.Do(req)
	// if err != nil {
	// 	log.Printf("Failed to make request: %v", err)
	// 	return err
	// }
	// log.Println("Request completed")

	return nil
}
