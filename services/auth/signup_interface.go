package auth

import (
	"gofr.dev/pkg/gofr"
	"github.com/Phase-R/Phase-R-Backend/db/models"
)

type User interface {
	CreateUser(ctx *gofr.Context, user *models.User) (*models.User, error)
	GetUser(ctx *gofr.Context, CUID string) (*models.User, error)
}
