package main

import (
	"flashquiz-server/internal/handlers/user"
	"flashquiz-server/internal/middlewares"
	"log"
	"net/http"
	"time"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/api/hello", user.Welcome)

	handler := middlewares.CORS(router)

	srv := http.Server{
		Addr: ":8080",
		Handler: handler,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 10*time.Second,
		IdleTimeout: 15*time.Second,
	}

	log.Println("Server running on port 8080")
	log.Fatal(srv.ListenAndServe())
}