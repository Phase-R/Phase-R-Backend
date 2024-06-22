package auth

import (
	"database/sql"
	"errors"
	"github.com/Phase-R/Phase-R-Backend/auth/tools"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/lib/pq"
	"github.com/nrednav/cuid2"
	"gofr.dev/pkg/gofr"
)

type user struct{}

func New() User {
	return user{}
}

func (u user) CreateUser(ctx *gofr.Context, user *models.User) (*models.User, error) {
	const uniqueViolation = "23505"

	id := cuid2.Generate()
	if id == "" {
		return nil, errors.New("CUID Generation failure.")
	}

	_, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO user (CUID, Username,Fname, Lname, Email, Password, Age, Access) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, user.Username, user.Fname, user.Lname, user.Email, tools.PwdSaltAndHash(user.Password), user.Age, user.Access)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == uniqueViolation {
				return nil, errors.New("entity already exists")
			}
		}
		return nil, errors.New("DB error")
	}
	return user, nil
}

func (u user) FetchUser(ctx *gofr.Context, CUID string) (*models.User, error) {
	var resp models.User

	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT CUID, Username, Fname, Lname, Email, Password, Age, Access FROM users WHERE CUID=$1", CUID).Scan(&resp.CUID, &resp.Username, &resp.Fname, &resp.Lname, &resp.Email, &resp.Password, &resp.Age, &resp.Access)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err
	}

	return &resp, nil
}
