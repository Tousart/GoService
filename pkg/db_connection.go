package pkg

import (
	"database/sql"
	"fmt"
)

func ConnectToDB(host *string, port *uint16, name *string, ssl *string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://user:password@%s:%d/%s?sslmode=%s", *host, *port, *name, *ssl)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверка: установлено ли соединение с бд
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
