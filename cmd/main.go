package main

import (
	"flashquiz-server/internal/middlewares"
	"flashquiz-server/internal/handlers"
	"flashquiz-server/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	err := database.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	router := http.NewServeMux()

	router.HandleFunc("/api", handlers.Welcome)
	router.HandleFunc("/api/register", handlers.Register)
	router.HandleFunc("/api/login", handlers.Login)
	router.HandleFunc("/api/user", handlers.UserHandler)

	handler := middlewares.CORS(middlewares.Recovery(middlewares.AuthMiddleware(router)))

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	srv := &http.Server{
		Addr: ":"+PORT,
		Handler: handler,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 10*time.Second,
		IdleTimeout: 15*time.Second,
	}

	log.Println("Connected to Postgres")
	log.Printf("Server running on port %s\n", PORT)
	log.Fatal(srv.ListenAndServe())
}