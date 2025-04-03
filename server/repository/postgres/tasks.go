package postgres

import (
	"database/sql"
	"errors"
	"httpServer/server/config"
	"httpServer/server/domain"
	"httpServer/server/repository"

	_ "github.com/lib/pq"
)

type Tasks struct {
	db *sql.DB
}

func NewTasksRepository(cfg config.Postgres) (*Tasks, error) {
	// connStr := fmt.Sprintf("postgres://user:password@%s:%d/%s?sslmode=%s", cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	return nil, err
	// }

	// // Проверка: установлено ли соединение с бд
	// if err := db.Ping(); err != nil {
	// 	return nil, err
	// }
	db, err := connectToDB(&cfg.Host, &cfg.Port, &cfg.DBName, &cfg.SSLMode)
	if err != nil {
		return nil, err
	}

	if err := applyMigrations(db); err != nil {
		return nil, err
	}

	return &Tasks{
		db: db,
	}, nil
}

// Возвращаем статус таски
func (rs *Tasks) GetStatus(taskId string) (string, error) {
	var value string
	err := rs.db.QueryRow("SELECT status FROM tasks WHERE task_id = $1", taskId).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", repository.ErrNotFound
	}

	return value, nil
}

// Возвращаем результат таски
func (rs *Tasks) GetResult(taskId string) (string, string, error) {
	var (
		stdout string
		stderr string
	)
	err := rs.db.QueryRow("SELECT stdout, stderr FROM tasks WHERE task_id = $1", taskId).Scan(&stdout, &stderr)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", repository.ErrNotFound
	}

	return stdout, stderr, nil
}

// Создаем статус и результат таски
func (rs *Tasks) PostTask(task *domain.Task) error {
	_, err := rs.db.Exec("INSERT INTO tasks (task_id, status, stdout, stderr) VALUES ($1, $2, $3, $4)", task.TaskId, task.Status, task.Stdout, task.Stderr)
	if err != nil {
		return err
	}

	// _, err := rs.db.Exec(`INSERT INTO tasks (task_id, status, stdout, stderr) VALUES ($1, $2, $3, $4)
	// ON CONFLICT (task_id) DO UPDATE SET
	// status = EXCLUDED.status,
	// stdout = EXCLUDED.stdout,
	// stderr = EXCLUDED.stderr`, task.TaskId, task.Status, task.Stdout, task.Stderr)
	// if err != nil {
	// 	return err
	// }

	return nil
}
