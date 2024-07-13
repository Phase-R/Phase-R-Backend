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
		c.JSON(405, gin.H{
			"error": "invalid email or password (email)",
		})
		return
	}
	// if user.Password != body.Password {
	// 	c.JSON(405, gin.H{
	// 		"error": "invalid email or password (password)",
	// 	})
	// 	return
	// }

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
			"error": "token generation error",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*30,"","", false, true)

	c.JSON(200, gin.H{
		"message": "login successful",
	})
}

func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User 
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(405, gin.H{"error": "Email not found, please sign up"})
		return
	}

	otp, hashedOTP, err := utils.GenerateOTP()
	if err!=nil {
		c.JSON(401, gin.H{"error":"error in hashing otp"})
	}

	user.otp= hashedOTP
	res := db.DB.Model(&models.User{}).Where("email = ?", body.Email).Update("otp", hashedOTP)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	err = sendEmailOTP(user.Email, otp)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send OTP to email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email","hashedOTP":hashedOTP,"otp":otp})
}