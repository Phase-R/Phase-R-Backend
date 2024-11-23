package main

import (
	"log"

	"github.com/Phase-R/Phase-R-Backend/services/auth/controllers"
	"github.com/Phase-R/Phase-R-Backend/services/auth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {								
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	
	config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"}
    config.AllowHeaders = []string{"Content-Type"}
	config.AllowCredentials = true
    r.Use(cors.New(config))

	db.Init()
	r.POST("/user/new", controllers.CreateUser)
	r.POST("/user/login", controllers.Login)
	r.GET("/verify", controllers.VerifyEmail)

	r.POST("/user/forgot-password", controllers.ForgotPassword)			
	r.POST("/user/reset-password", controllers.ResetPassword)
	// r.GET("/user/google/signin",controllers.BeginGoogleAuth)
	// r.GET("user/google/callback", controllers.OAuthCallback)
	r.Run()
}