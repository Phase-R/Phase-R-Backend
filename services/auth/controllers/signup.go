package controllers

// package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/Phase-R/Phase-R-Backend/auth/utils"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nrednav/cuid2"
)

func getJWTSecretKey() string {
	return os.Getenv("JWT_SECRET")
}

func generateVerificationToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"ttl":    time.Now().Add(time.Minute * 5).Unix(),
	})

	jwtSecret := getJWTSecretKey()
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	fmt.Println("token: ", tokenString)

	return tokenString, nil
}

func CreateUser(ctx *gin.Context) {
	const uniqueViolation = "23505"

	var newUser models.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.Where("email = ?", newUser.Email).First(&newUser).Error
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
	}

	id := cuid2.Generate()
	if id == "" {
		return
	}

	newUser.ID = id

	hash, err := utils.PwdSaltAndHash(newUser.Password)
	if err != nil {
		log.Fatal("could not hash password", err)
	}

	newUser.Password = hash

	res := db.DB.Create(&newUser)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send verification email
	token, err := generateVerificationToken(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// send email
	err = SendVerificationEmail(newUser.Email, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Verification email sent."})
}

func FetchUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	res := db.DB.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedUser models.User

	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.DB.Model(&models.User{}).Where("id = ?", id).Updates(updatedUser)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	res := db.DB.Delete(&models.User{}, id)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
