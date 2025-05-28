package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func Verify(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func GenerateJWT(id int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"exp": time.Now().Add(24*time.Hour).Unix(),
	}

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, ok in prod")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(os.Getenv("JWT_KEY"))
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

func (user User) LoginValidate() error {
	if user.Email == "" {
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