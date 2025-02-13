package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (r Repository) getUserIDByUsername(ctx context.Context, tx pgx.Tx, username string) (int, error) {
	var userID int
	err := tx.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r Repository) getUserBalance(ctx context.Context, tx pgx.Tx, userID int) (int, error) {
	var balance int
	err := tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
