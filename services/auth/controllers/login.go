package controllers

import (
	"os"
	"time"
	"github.com/Phase-R/Phase-R-Backend/services/auth/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/auth/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"net/http"
	"errors"
	"github.com/alexedwards/argon2id"
)

func ParseToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func Login(c *gin.Context) {
	tokenString, err := c.Cookie("Auth")
	if err == nil {
		token, err := ParseToken(tokenString)
		if err == nil && token.Valid {
			c.JSON(201, gin.H{
				"message": "already logged in",
			})
			return
		}
	}

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
	if result.Error != nil{
		c.JSON(404, gin.H{
			"error": "invalid email or password (email)",
		})
		return
	}
	if !user.Verified {
		c.JSON(405, gin.H{
			"error": "email not verified",
		})
		return
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil || !match {
		c.JSON(404, gin.H{
			"error": "invalid email or password compare",
		})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	
	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(401, gin.H{
			"error": "token ge		neration error",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*30,"","", false, false)

	c.JSON(200, gin.H{
		"message": "login successful",
	})
}

func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User

	// Find the user by email
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(405, gin.H{"error": "Email not found, please sign up"})
		return
	}

	// Generate OTP and hashed OTP
	otp, hashedOTP, err := utils.GenerateOTP()
	if err != nil {
		c.JSON(401, gin.H{"error": "Error generating OTP"})
		return
	}

	// Update the user's OTP field
	user.OTP = hashedOTP

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update OTP in the database"})
		return
	}

	// Send the OTP email
	err = sendEmailOTP(user.Email, otp)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send OTP to email"})
		return
	}
	c.JSON(200, gin.H{
		"message": "login successful",
	})
}

func ResetPassword(c *gin.Context) {
	var body struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		Password 	string `json:"password"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	var user models.User

	// Find the user by email
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Verify the OTP
	match, err := argon2id.ComparePasswordAndHash(body.OTP, user.OTP)
	if err != nil || !match {
		c.JSON(404, gin.H{
			"error": "invalid email or password compare",
		})
		return
	}

	// Hash the new password
	hashedPassword, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update the user's password
	user.Password = hashedPassword
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}

	// Generate a new JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(401, gin.H{"error": "Token generation error"})
		return
	}

	// Set the new JWT token as a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}