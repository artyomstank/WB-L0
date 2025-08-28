package repo

import (
	"database/sql"
)

// PostgresRepo содержит *sql.DB и методы для работы с таблицами
type PostgresRepo struct {
	DB *sql.DB
}

// Конструктор PostgresRepo
func NewRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

// Закрытие соединения
func (r *PostgresRepo) Close() error {
	if r.DB != nil {
		return r.DB.Close()
	}
	return nil
}
