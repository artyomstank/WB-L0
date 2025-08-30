package db

import (
	"L0-wb/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Подключение к Postgres
func NewDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.Postgres.GetDBConnStr())
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}
	fmt.Println("DB connected")
	return db
}
