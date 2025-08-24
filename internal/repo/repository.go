package repo

import "database/sql"

type PostgresRepo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) PostgresRepo {
	return PostgresRepo{
		DB: db,
	}
}
