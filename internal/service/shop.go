package service

import (
	"context"
)

func (s Service) BuyItem(ctx context.Context, userID int, itemName string) error {
	return s.repo.BuyItem(ctx, userID, itemName)
}
