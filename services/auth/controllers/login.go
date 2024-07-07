package controllers

import (
	// "fmt"
	"log"
	"os"
	"time"
	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/auth/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
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
	if result.Error != nil{
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
	log.Println(user.Password,body.Password)
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
	c.SetCookie("Auth", token, 3600*24*30,"","", false, true)

	c.JSON(200, gin.H{
		"message": "login successful",
	})
}