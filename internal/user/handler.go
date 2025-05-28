package user

import (
	"database/sql"
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
		ErrorResponse(w, "Invalid JSON Format")
		return
	}

	err = u.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorResponse(w, "name, email & password must be provided")
		return
	}

	var userExists bool
	if err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", u.Email).Scan(&userExists); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if userExists {
		http.Error(w, "User already exists, Try logging in", http.StatusNotAcceptable)
		ErrorResponse(w, "User already exists, Try logging in")
		return
	}

	hashedPassword, err := Hash(u.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return		
	}

	_, err = db.DB.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", u.Name, u.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return		
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON Format", http.StatusBadRequest)
		ErrorResponse(w, "Invalid JSON Format")
		return
	}

	if err := u.LoginValidate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorResponse(w, "email & password must be provided")
		return
	}

	var user User
	if err := db.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", u.Email).Scan(&user.Id, &user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		if err == sql.ErrNoRows {
			ErrorResponse(w, "User dosen't exist")
			return
		}
		ErrorResponse(w, "Failed to check user")
		return
	}

	err = Verify(user.Password, u.Password)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		ErrorResponse(w, "Invalid Crendentials")
		return
	}

	token, err := GenerateJWT(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ErrorResponse(w, "Failed to generate jwt")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"message": "Login Successful",
	})
}