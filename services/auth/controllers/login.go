package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB



// func sync_database(){
// 	db.AutoMigrate(&models.User{})
// }




func connect_to_database()*gorm.DB{
	dbHost := os.Getenv("dbHost")
    dbPort := os.Getenv("dbPort")
    dbName := os.Getenv("dbName")
    dbUser := os.Getenv("dbUser")
    dbPassword := os.Getenv("dbPassword")
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbName, dbPort)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{});if err!=nil{
		log.Fatal("could not connect to db")
	}
	return db
}

func init(){
	err:=godotenv.Load("configs/.env")

	if err!=nil{
		log.Fatal("Error loading .env file")
	}
}

func login(c *gin.Context) {
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
	query:="SELECT * FROM users WHERE email = ?"
	result := db.Raw(query, body.Email).Scan(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(405, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	if user.Password != body.Password {
		c.JSON(405, gin.H{
			"error": "invalid email or password",
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




// func main(){
// 	r:=gin.Default()
// 	// sync_database()
// 	r.GET("/ping",func(ctx *gin.Context) {
// 		ctx.JSON(200,gin.H{
// 			"message":"pong",
// 		})
// 	})
// 	db=connect_to_database()
// 	r.POST("/login",login)

// 	r.Run()
// }