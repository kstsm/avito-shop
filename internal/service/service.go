package service

import (
	"context"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/repository"
)

type ServiceI interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error)
	BuyItem(ctx context.Context, userID int, item string) error
	SendCoins(ctx context.Context, amount, senderID int, username string) error
}

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) *Service {
	return &Service{
		repo: repo,
	}
}
