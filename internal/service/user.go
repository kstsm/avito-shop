package service

import (
	"context"
	"github.com/kstsm/avito-shop/api/rest/models"
)

func (s Service) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	return s.repo.GetUserInfo(ctx, userID)
}

func (s Service) SendCoins(ctx context.Context, amount, senderID int, username string) error {
	return s.repo.SendCoins(ctx, senderID, amount, username)
}
