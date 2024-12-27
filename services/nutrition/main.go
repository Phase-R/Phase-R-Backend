package main

import (
	// "log"
	
	"github.com/Phase-R/Phase-R-Backend/services/nutrition/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.POST("/monthly_diet_gen", controllers.Monthly_Diet_Gen)
	r.Run()
}