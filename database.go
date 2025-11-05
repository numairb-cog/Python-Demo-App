package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(config *Config) error {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDB,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		log.Printf("Warning: Could not connect to database: %v", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
