package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var DB *sql.DB

func InitDB() error {
	var err error
	// Still in development
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=enable", "test", "test", "test", "test")

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return errors.New("failed to initialize database connection")
	}
	defer DB.Close()

	// Check the connection status
	if err := DB.Ping(); err != nil {
		return errors.New("failed to connect to database")
	}

	initTables(DB)

	return nil
}

func initTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY, 
			name VARCHAR(100) NOT NULL, 
			email VARCHAR(100) UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return errors.New("database error")
	}

	return nil
}