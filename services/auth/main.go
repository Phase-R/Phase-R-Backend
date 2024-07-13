package main

import (
	"log"

	"github.com/Phase-R/Phase-R-Backend/auth/controllers"
	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.Use(cors.Default())
	db.Init()
	r.POST("/user/new", controllers.CreateUser)
	r.POST("/user/login", controllers.Login)
	r.GET("/verify", controllers.VerifyEmail)
	r.Run()

}
