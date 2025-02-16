package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/api/rest/models"
)

type RepositoryI interface {
	CreateUser(ctx context.Context, username, hashedPassword string) (int, error)
	GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error)
	BuyItem(ctx context.Context, userID int, itemName string) error
	SendCoins(ctx context.Context, senderID, amount int, username string) error
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) RepositoryI {
	return &Repository{
		conn: conn,
	}
}
