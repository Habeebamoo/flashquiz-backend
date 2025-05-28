package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() error {
	var err error

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, ok in prod")
	}

	connStr := os.Getenv("DATABASE_URL")

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return errors.New("failed to initialize database connection")
	}

	// Check the connection status
	if err := DB.Ping(); err != nil {
		return errors.New("failed to connect to database")
	}

	return nil
}