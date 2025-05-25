package main

import (
	"flashquiz-server/internal/middlewares"
	"flashquiz-server/internal/user"
	"flashquiz-server/pkg/db"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/api/", user.Welcome)
	router.HandleFunc("/api/register", user.Register)

	handler := middlewares.CORS(middlewares.Recovery(router))

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
	log.Println("Server running on port 8080")
	log.Fatal(srv.ListenAndServe())
}