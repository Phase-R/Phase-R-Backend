package auth

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"log"
	"net/http"
	"os"

	"github.com/Phase-R/Phase-R-Backend/auth/tools"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

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

//func main() {
//	r := gin.Default()
//	r.GET("/ping", func(ctx *gin.Context) {
//		ctx.JSON(http.StatusOK, gin.H{"yohoo": "pong"})
//	})
//	db = dbConn()
//	r.POST("/user/new", CreateUser)
//	r.GET("/user/fetch", FetchUser)
//	r.Run()
//}

//var db *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func dbConn() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=enabled", host, user, passwd, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failure : %v", err)
	}
	return db
}
