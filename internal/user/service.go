package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (user User) Validate() error {
	if user.Name == "" {
		return errors.New("name is not defined")
	} else if user.Email == "" {
		return errors.New("email is not defined")
	} else if user.Password == "" {
		return errors.New("password is not defined")
	} else {
		return nil
	}
}

func ErrorResponse(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}