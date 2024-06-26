package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	// "github.com/gofiber/fiber/v2"
	// "net/http"
	// "github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gofr.dev/pkg/gofr"
	// resTypes "gofr.dev/pkg/gofr/http/response"
)

var db *sql.DB

func init_secret() string {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET_KEY")
	return secretKey
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(ctx *gofr.Context) (interface{}, error) {
	var loginRequest LoginRequest
	if err := ctx.Bind(&loginRequest); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var Email, password string
	query := `SELECT email, password FROM "User" WHERE email = $1`
	row := db.QueryRow(query, loginRequest.Email)
	if err := row.Scan(&Email, &password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}
		ctx.Logger.Errorf("error querying the database: %v", err)
		return nil, err
	}
	if password != loginRequest.Password {
		return nil, fmt.Errorf("invalid email or password")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(init_secret()))
	if err != nil {
		return nil, err
	}
	// TODO : COOKIE SESSION MANAGEMENT

	// expiration := time.Now().Add(24 * time.Hour)
	// cookie := http.Cookie{
	//     Name:     "jwt",
	//     Value:    token,
	//     Expires:  expiration,
	//     HttpOnly: true,
	// }
	// http.SetCookie(gofr.Responder.Respond(cookie), &cookie)
	// http.SetCookie(gofr.Responder, &cookie)
	return token, nil
}

func main() {

}
