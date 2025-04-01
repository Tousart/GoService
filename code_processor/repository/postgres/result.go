package postgres

import (
	"database/sql"
	"fmt"
	"httpServer/code_processor/config"
	"httpServer/code_processor/domain"
)

type ResultRepository struct {
	db *sql.DB
}

func NewResultRepository(cfg config.Postgres) (*ResultRepository, error) {
	connStr := fmt.Sprintf("postgres://user:password@%s:%d/%s?sslmode=%s", cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &ResultRepository{
		db: db,
	}, nil
}

func (rs *ResultRepository) SendResult(result *domain.Result) error {
	_, err := rs.db.Exec(`INSERT INTO tasks (task_id, status, stdout, stderr) VALUES ($1, $2, $3, $4)
	ON CONFLICT (task_id) DO UPDATE SET
	status = EXCLUDED.status,
	stdout = EXCLUDED.stdout,
	stderr = EXCLUDED.stderr`, result.TaskId, result.Status, result.Stdout, result.Stderr)
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
