package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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

func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateJWT(id string) (string, error) {
	claims := jwt.MapClaims{
		"userId": id,
		"exp": time.Now().Add(24*time.Hour).Unix(),
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no .env file, ok in prod")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func ErrorResponse(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

func SendVerification(userEmail, userName, token string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "flashquizweb@gmail.com")
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Verify your FlashQuiz account.")

	body := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<h2 style="color: #1a73e8">Welcome to FlashQuiz</h2>
				<p>Dear %s</p>

				<p>Thank you for creating an account on <strong>FlashQuiz</strong>, the ultimate destination for trivia and fun.</p>

				<p>To complete your registration and activate your account, please confirm your email address by clicking the button below:</p>

				<p><a href="https://flashquizweb.netlify.app/verify?token=%s" style="color: white; background-color: #1a73e8; padding: 20px; display: block; text-align: center; font-weight: bold; font-size: 1.2em;">Verify My Email</a></p>

				<p>This verification grants you access to all our website features and quizzes, it also helps us secure your account and ensure only you have access to it.</p>

				<p>If you didn't sign up for FlashQuiz, you can safely ignore this email</p>
				<br>
				<p>Thank you for choosing FlashQuiz.</p>

				<br>
				<p>Best regards,<br>FlashQuiz Team</p>
			</body>
		</html>
	`, userName, token)

	m.SetBody("text/html", body)

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no .env file, ok in prod")
	}

	d := gomail.NewDialer("smtp.gmail.com", 465, "flashquizweb@gmail.com", os.Getenv("APP_PASSWORD"))
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
