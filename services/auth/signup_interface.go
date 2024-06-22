package auth

import (
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"gofr.dev/pkg/gofr"
)

type User interface {
	CreateUser(ctx *gofr.Context, user *models.User) (*models.User, error)
	GetUser(ctx *gofr.Context, CUID string) (*models.User, error)
}
