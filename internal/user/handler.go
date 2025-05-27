package user

import (
	"encoding/json"
	"flashquiz-server/pkg/db"
	"fmt"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello from FlashQuiz Backend Server")
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON Format", http.StatusBadRequest)
		return
	}

	err = u.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userExists bool
	if err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", u.Email).Scan(&userExists); err != nil {
		http.Error(w, "Internal Server Error (userExists)", http.StatusInternalServerError)
		return
	}

	if userExists {
		http.Error(w, "User already exists, Try logging in", http.StatusNotAcceptable)
		return
	}

	hashedPassword, err := Hash(u.Password)
	if err != nil {
		http.Error(w, "Internal Server Error (hashing password)", http.StatusInternalServerError)
		return		
	}

	_, err = db.DB.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", u.Name, u.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Internal Server Error (inserting database)", http.StatusInternalServerError)
		return		
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}