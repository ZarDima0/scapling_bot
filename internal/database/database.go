package database

import (
	"database/sql"
	"fmt"
	"github.com/ZarDima0/scapling_bot/internal/config"
	_ "github.com/lib/pq"
	"log"
)

func NewDB(cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.SslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
