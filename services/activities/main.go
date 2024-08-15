package main

import (
	"github.com/Phase-R/Phase-R-Backend/activities/controllers"
	"github.com/Phase-R/Phase-R-Backend/activities/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	db.Init()
	//drill endpoints
	r.GET("/get_drill/:id", controllers.GetDrill)
	r.GET("/get_drills_by_type/:type", controllers.GetDrillsByType)
	r.PUT("/update_drill/:id", controllers.UpdateDrill)
	r.DELETE("/delete_drill/:id", controllers.DeleteDrill)
	// Activity endpoints
	r.POST("/create_activity", controllers.CreateActivity)
	r.GET("/get_activity/:id", controllers.GetActivity)
	r.PUT("/update_activity/:id", controllers.UpdateActivity)
	r.DELETE("/delete_activity/:id", controllers.DeleteActivity)
	// ActivityType endpoints
	r.POST("/create_activity_type", controllers.CreateActType)
	r.GET("/get_activity_type/:id", controllers.GetActType)
	r.PUT("/update_activity_type/:id", controllers.UpdateActType)
	r.DELETE("/delete_activity_type/:id", controllers.DeleteActType)
	// User progress endpoints
	r.POST("/user_progress",controllers.ProgressController)
	r.POST("/get_user_progress",controllers.GetUserProgress)
	r.Run()
}
