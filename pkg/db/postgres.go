package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() error {
	var err error

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file, ok in prod")
	}
	// Still in development
	connStr := os.Getenv("DATABASE_URL")

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return errors.New("failed to initialize database connection")
	}

	// Check the connection status
	if err := DB.Ping(); err != nil {
		return errors.New("failed to connect to database")
	}

	if err := initTables(DB); err != nil {
		return err
	}

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