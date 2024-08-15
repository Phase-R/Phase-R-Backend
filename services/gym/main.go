package main

import (
	"github.com/Phase-R/Phase-R-Backend/gym/controllers"
	"github.com/Phase-R/Phase-R-Backend/gym/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	db.Init()
	r.GET("/get_all_workouts",controllers.GetAllGymData)
	r.Run(":8082")
}