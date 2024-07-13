package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Phase-R/Phase-R-Backend/services/auth/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	gomail "gopkg.in/gomail.v2"
)

func SendVerificationEmail(emailTo string, token string) error {
	password := os.Getenv("MAIL_PASS")
	from := os.Getenv("EMAIL_FROM")

	header, err := os.ReadFile("./controllers/email_templates/verify_email_header.html")
	if err != nil {
		fmt.Println("Error reading header file: ", err)
		return err
	}

	footer, err := os.ReadFile("./controllers/email_templates/verify_email_footer.html")
	if err != nil {
		fmt.Println("Error reading footer file: ", err)
		return err
	}

	body := string(header) + fmt.Sprintf("<a href='http://localhost:8080/verify?token=%s'>Verify Email</a>", token) + string(footer)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", "Verify your email")
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

	return nil
}

func VerifyEmail(ctx *gin.Context) {
	tokenString := ctx.Query("token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getJWTSecretKey()), nil
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JWT token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "JWT Claims failed!"})
		return
	}

	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}

	// fmt.Println("Claims user id: ", claims["userID"])

	var user models.User
	db.DB.Raw("SELECT * FROM users WHERE id = ?", claims["userID"]).Scan(&user)

	if user.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already verified"})
		return
	}

	user.Verified = true
	result := db.DB.Save(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}


	GenerateTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	GenerateTokenstring, err := GenerateTokenClaims.SignedString([]byte(getJWTSecretKey()))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate new JWT token"})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("verified", GenerateTokenstring, 3600*24*30,"","", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Email successfully verified"})
}

func sendEmailOTP(emailTo string, otp int) error {
	password := os.Getenv("MAIL_PASS")
	from := os.Getenv("EMAIL_FROM")

	header, err := os.ReadFile("./controllers/email_templates/verify_email_header.html")
	if err != nil {
		fmt.Println("Error reading header file: ", err)
		return err
	}

	footer, err := os.ReadFile("./controllers/email_templates/verify_email_footer.html")
	if err != nil {
		fmt.Println("Error reading footer file: ", err)
		return err
	}

	body := string(header) + fmt.Sprintf("<p>daddy likes you, heres his phone number <3:</p><h1>%d</h1>", otp) + string(footer)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", "Forget password OTP")
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}