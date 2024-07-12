package controllers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/Phase-R/Phase-R-Backend/auth/utils"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	gomail "gopkg.in/gomail.v2"
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
	if result.Error != nil {
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
	c.SetCookie("Auth", token, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{
		"message": "login successful",
	})
}

func ForgotPassword(c *gin.Context) {
    var body struct {
        Email string `json:"email"`
    }

    if c.Bind(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    // Generate OTP
    otp, err := GenerateOTP(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Generate a new JWT token
    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
        Issuer:    body.Email,
        ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Token valid for 1 hour
    })

    token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    // Send OTP via email
    sendEmail(body.Email, otp)

    c.JSON(http.StatusOK, gin.H{"message": "password reset OTP sent", "token": token})
}


// ResetPassword function
func ResetPassword(c *gin.Context) {
    var body struct {
        Email       string `json:"email"`
        OTP         string `json:"otp"`
        NewPassword string `json:"new_password"`
    }

    if c.Bind(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    var user models.User
    result := db.DB.Where("email = ?", body.Email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Compare the provided OTP with the hashed OTP stored in the database
    match, err := utils.ComparePasswords(user.OTP, body.OTP)
    if err != nil || !match {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
        return
    }

    // Hash the new password
    hash, err := utils.PwdSaltAndHash(body.NewPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    // Update the user's password and clear the OTP
    user.Password = hash
    user.OTP = "" // Clear the OTP after successful reset
    db.DB.Save(&user)

    // Generate a new JWT token
    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
        Issuer:    user.Email,
        ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Token valid for 1 hour
    })

    token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "password reset successful", "token": token})
}