package main

import (
	"github.com/Phase-R/Phase-R-Backend/auth/controllers"
	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	db.Init()
	r.POST("/user/new",controllers.CreateUser)
	r.POST("/user/login",controllers.Login)
	r.Run()
}
