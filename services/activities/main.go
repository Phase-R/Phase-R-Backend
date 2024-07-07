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
	db := database.InitDB()
	// Drill endpoints
	r.POST("/create_drill", controllers.CreateDrill(db))
	r.GET("/get_drill/:id", controllers.GetDrill(db))
	r.GET("/get_drills_by_type/:type", controllers.GetDrillsByType(db))
	r.PUT("/update_drill/:id", controllers.UpdateDrill(db))
	r.DELETE("/delete_drill/:id", controllers.DeleteDrill(db))
	// Activity endpoints
	r.POST("/create_activity", controllers.CreateActivity(db))
	r.GET("/get_activity/:id", controllers.GetActivity(db))
	r.PUT("/update_activity/:id", controllers.UpdateActivity(db))
	r.DELETE("/delete_activity/:id", controllers.DeleteActivity(db))
	r.Run()
}
