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
	r.GET("/get_drill/:id", controllers.GetDrill)
	r.GET("/get_drills_by_type/:type", controllers.GetDrillsByType)
	r.PUT("/update_drill/:id", controllers.UpdateDrill)
	r.DELETE("/delete_drill/:id", controllers.DeleteDrill)
	// Activity endpoints
	r.POST("/create_activity", controllers.CreateActivity)
	r.GET("/get_activity/:id", controllers.GetActivity)
	r.PUT("/update_activity/:id", controllers.UpdateActivity)
	r.DELETE("/delete_activity/:id", controllers.DeleteActivity)
	r.Run()
}
