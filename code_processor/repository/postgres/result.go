package postgres

import (
	"database/sql"
	"httpServer/code_processor/config"
	"httpServer/code_processor/domain"
	"httpServer/pkg"
)

type ResultRepository struct {
	db *sql.DB
}

func NewResultRepository(cfg config.Postgres) (*ResultRepository, error) {
	db, err := pkg.ConnectToDB(&cfg.Host, &cfg.Port, &cfg.DBName, &cfg.SSLMode)
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

	return nil
}
