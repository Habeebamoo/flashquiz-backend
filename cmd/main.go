package main

import (
	"flashquiz-server/internal/handlers/user"
	"flashquiz-server/internal/middlewares"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/api/hello", user.Welcome)

	handler := middlewares.CORS(router)

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

	log.Println("Server running on port 8080")
	log.Fatal(srv.ListenAndServe())
}