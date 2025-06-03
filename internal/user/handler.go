package user

import (
	"database/sql"
	"encoding/json"
	"flashquiz-server/internal/middlewares"
	"flashquiz-server/pkg/db"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorResponse(w, "Method Not Allowed")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "This is FlashQuiz Server",
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorResponse(w, "Method Not Allowed")
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse(w, "Invalid JSON Format")
		return
	}

	err = u.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse(w, err.Error())
		return
	}

	var userExists bool
	if err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", u.Email).Scan(&userExists); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorResponse(w, "Internal Server Error")
		return
	}

	if userExists {
		w.WriteHeader(http.StatusConflict)
		ErrorResponse(w, "User already exists, Try logging in")
		return
	}

	hashedPassword, err := Hash(u.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorResponse(w, "Internal Server Error")
		return		
	}

	_, err = db.DB.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", u.Name, u.Email, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorResponse(w, "Internal Server Error")
		return		
	}

	go func() {
		err := SendVerification(u.Email, u.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorResponse(w, "Internal Server Error")
		}
	}()

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorResponse(w, "Method Not Allowed")
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse(w, "Invalid JSON Format")
		return
	}

	if err := u.LoginValidate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse(w, "Email and password must be provided")
		return
	}

	var user User
	if err := db.DB.QueryRow("SELECT user_id, password FROM users WHERE email = $1", u.Email).Scan(&user.UserId, &user.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err == sql.ErrNoRows {
			ErrorResponse(w, "User not found")
			return
		}
		ErrorResponse(w, "Internal Server Error")
		return
	}

	err = Verify(user.Password, u.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ErrorResponse(w, "Invalid Crendentials")
		return
	}

	token, err := GenerateJWT(user.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorResponse(w, "Method Not Allowed")
		return
	}

	userId, ok := r.Context().Value(middlewares.UserIdKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		ErrorResponse(w, "Unauthorized Access")
		return
	}

	var user UserResponse
	err := db.DB.QueryRow("SELECT user_id, name, email, isVerified FROM users WHERE id = $1", userId).Scan(&user.UserId, &user.Name, &user.Email, &user.IsVerified)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorResponse(w, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]UserResponse{
		"data": user,
	})
}