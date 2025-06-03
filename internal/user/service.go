package user

import (
	"encoding/json"
	"errors"
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

func GenerateJWT(id string) (string, error) {
	claims := jwt.MapClaims{
		"userId": id,
		"exp": time.Now().Add(24*time.Hour).Unix(),
	}

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, ok in prod")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
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

func SendVerification(userEmail, userName string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "habeebamoo08@gmail.com")
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Verify your FlashQuiz account.")

	body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello %s</p>
				<p>Thank you for signing up on <strong>FlashQuiz</strong>! To activate your account, pease verify your email address by clicking on the link below:</p>
				<p><a href="%s" style="color: #1a73e8;">Verify My Email</a></p>
				<p>If you didn't create a FlashQuiz account, you can safely ignore this email.</p>
				<br>
				<p>Best regards,<br>FlashQuiz Team</p>
			</body>
		</html>
	`, userName, "https://flashquizweb.netlify.app")

	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 465, "habeebamoo08@gmail.com", "exphqvkpdzrbhrdp")
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		return errors.New("failed to send verification email")
	}

	return nil
}