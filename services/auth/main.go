package main

import (
	"log"
	"os"

	"github.com/Phase-R/Phase-R-Backend/services/auth/controllers"
	"github.com/Phase-R/Phase-R-Backend/services/auth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err) // Log detailed error message
	}

	// Initialize the Gin router
	r := gin.Default()

	// Initialize session middleware in controllers
	controllers.InitializeSessionMiddleware(r)

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Update this based on your frontend URL
	config.AllowHeaders = []string{"Content-Type"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Initialize database connection
	err = db.Init() // Handle database initialization error
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Initialize Google OAuth
	controllers.InitializeOAuthGoogle()

	// Routes
	r.POST("/user/new", controllers.CreateUser)
	r.POST("/user/login", controllers.Login)
	r.GET("/verify", controllers.VerifyEmail)
	r.POST("/user/forgot-password", controllers.ForgotPassword)
	r.POST("/user/reset-password", controllers.ResetPassword)

	// Google OAuth Routes
	r.GET("/user/google/signin", controllers.HandleGoogleLogin)
	r.GET("/user/google/auth", controllers.CallBackFromGoogle)

	r.Run()
}