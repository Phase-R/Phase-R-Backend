package main

import (
	"log"

	"github.com/Phase-R/Phase-R-Backend/activities/controllers"
	"github.com/Phase-R/Phase-R-Backend/db/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	Init()
	r := gin.Default()
	_ = database.InitDB()
	r.POST("/create_drill", controllers.CreateDrill)
	r.GET("/get_drill/:id", controllers.GetDrill)
	r.GET("/get_drills_by_type/:type", controllers.GetDrillsByType)
	r.PUT("/update_drill/:id", controllers.UpdateDrill)
	r.DELETE("/delete_drill/:id", controllers.DeleteDrill)
	r.Run()
}
