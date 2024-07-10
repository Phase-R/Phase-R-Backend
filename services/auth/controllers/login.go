package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/Phase-R/Phase-R-Backend/auth/utils"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/gomail.v2"
)

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(402, gin.H{
			"error": "fail to read the body",
		})
		return
	}
	var user models.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(405, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	// if user.Password != body.Password {
	// 	c.JSON(405, gin.H{
	// 		"error": "invalid email or password",
	// 	})
	// 	return
	// }
	match, err := utils.ComparePasswords(user.Password, body.Password)
	log.Println(user.Password, body.Password)
	if err != nil || !match {
		c.JSON(404, gin.H{
			"error": "invalid email or password compare",
		})
		return
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(401, gin.H{
			"error": "token generation error",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{
		"message": "login successful",
	})
}

// ForgotPassword function
// ForgotPassword function
func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	// Generate OTP for user
	otp, err := GenerateOTP(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate OTP"})
		return
	}

	// Send OTP via email
	sendEmail(user.Email, otp)

	c.JSON(http.StatusOK, gin.H{"message": "password reset OTP sent"})
}

// ResetPassword function
// ResetPassword function
func ResetPassword(c *gin.Context) {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Parse and verify the token
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("RESET_PASSWORD_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["iss"].(string)
		var user models.User
		result := db.DB.Where("email = ?", email).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Verify OTP
		// if !utils.VerifyPassword(user.OTP, body.NewPassword) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OTP"})
		// 	return
		// }

		// Hash the new password
		hashedPassword, err := utils.PwdSaltAndHash(body.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}

		// Update user's password
		user.Password = hashedPassword
		db.DB.Save(&user)

		c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
	}
}

func sendEmail(to string, token string) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("MAIL_PASS")
	fmt.Println(from)
	fmt.Println(password)
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("<a href='http://localhost:8080/user/reset-password-test'>Click here</a>", token))

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("could not send email: ", err)
	}
}
