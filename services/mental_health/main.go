package main

import (
	"log"

	"github.com/Phase-R/Phase-R-Backend/services/mental_health/controllers"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file!")
	}

	r := gin.Default()
	config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"}
    config.AllowHeaders = []string{"Content-Type"}
	config.AllowCredentials = true
    r.Use(cors.New(config))

	db.Init()
	r.GET("/fetch_questions", controllers.FetchQuestionSet)
	r.POST("/evaluate_answers", controllers.ScoreEvaluation)
	r.Run()
}