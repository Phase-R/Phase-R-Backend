package main

import (
	"log"

	mentalHealthControllers "github.com/Phase-R/Phase-R-Backend/services/mental_health/controllers"
	mentalHealthDB "github.com/Phase-R/Phase-R-Backend/services/mental_health/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file!")
	}

	// Create a new Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Handle OPTIONS preflight requests
r.OPTIONS("/*path", func(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Status(204)
})

	// Initialize databases
	mentalHealthDB.Init()

	// Define routes for mental health services
	r.POST("/add_questions", mentalHealthControllers.AddQuestionSet)
	r.GET("/fetch_questions", mentalHealthControllers.FetchQuestionSet)
	r.POST("/evaluate_answers", mentalHealthControllers.ScoreEvaluation)
	r.POST("/post_thoughts", mentalHealthControllers.PostThoughts)
	// Start the Gin server
	r.Run()
}
