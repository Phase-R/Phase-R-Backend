package auth

import (
	"github.com/gin-gonic/gin"
)

type User interface {
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}
