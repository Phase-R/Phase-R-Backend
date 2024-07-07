//package auth

package main

import (
	"github.com/Phase-R/Phase-R-Backend/auth/tools"
	"gorm.io/gorm"

	"github.com/joho/godotenv"

	"github.com/Phase-R/Phase-R-Backend/db/database"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"github.com/nrednav/cuid2"
	"log"
	"net/http"
)

var db *gorm.DB

func CreateUser(ctx *gin.Context) {
	const uniqueViolation = "23505"

	var newUser models.User

	id := cuid2.Generate()
	if id == "" {
		return
	}

	newUser.ID = id

	hash, err := tools.PwdSaltAndHash(newUser.Password)
	if err != nil {
		log.Fatal("could not hash password", err)
	}

	newUser.Password = hash

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.Create(&newUser)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"yohoo": "new user created."})
}

func FetchUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	res := db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user)
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

	res := db.Model(&models.User{}).Where("id = ?", id).Updates(updatedUser)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	res := db.Delete(&models.User{}, id)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	Init()
	r := gin.Default()
	db = database.InitDB()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/user/new", CreateUser)
	r.GET("/user/fetch", FetchUser)
	err := r.Run(":5432")
	if err != nil {
		log.Fatal(err)
	}
}
