package auth

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Phase-R/Phase-R-Backend/auth/tools"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/lib/pq"
	"github.com/nrednav/cuid2"
	"gofr.dev/pkg/gofr"
)

func (user *models.User) CreateUser(ctx *gofr.Context, newuser *models.User) (*models.User, error) {
	const uniqueViolation = "23505"

	id := cuid2.Generate()
	if id == "" {
		return nil, errors.New("CUID Generation failure.")
	}

	hash, err := tools.PwdSaltAndHash(user.Password)
	if err != nil {
		log.Fatal("could not hash password", err)
	}

	_, err = ctx.SQL.ExecContext(ctx,
		"INSERT INTO user (id, username,fname, lname, email, password, age, access) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		id, user.Username, user.Fname, user.Lname, user.Email, hash, user.Age, user.Access)

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

func (user *models.User) FetchUser(ctx *gofr.Context, CUID string) (*models.User, error) {
	var resp models.User

	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, username, fname, lname, email, password, age, access FROM users WHERE id=?", CUID).Scan(&resp.ID, &resp.Username, &resp.Fname, &resp.Lname, &resp.Email, &resp.Password, &resp.Age, &resp.Access)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err
	}

	return &resp, nil
}
