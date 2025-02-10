package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/api/rest/models"
)

type RepositoryI interface {
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	CreateUser(ctx context.Context, username, password string) error
}

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) RepositoryI {
	return &Repository{
		conn: conn,
	}
}
