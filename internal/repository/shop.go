package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/internal/apperrors"
)

func (r Repository) BuyItem(ctx context.Context, userID int, itemName string) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		slog.Error("Ошибка при начале транзакции", "error", err)
		return fmt.Errorf("r.conn.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	var itemID, price int
	err = tx.QueryRow(ctx, QueryGetItem, itemName).Scan(&itemID, &price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Товар не найден", "item", itemName)
			return fmt.Errorf("%w: %s", apperrors.ErrItemNotFound, itemName)
		}
		slog.Error("Ошибка при получении товара", "error", err)
		return fmt.Errorf("tx.QueryRow: %w", err)
	}

	balance, err := r.getUserBalance(ctx, tx, userID)
	if err != nil {
		slog.Error("Ошибка при получении баланса", "userID", userID, "error", err)
		return fmt.Errorf("r.getUserBalance: %w", err)
	}

	if balance < price {
		slog.Warn("Недостаточно средств", "userID", userID, "balance", balance, "price", price)
		return apperrors.ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, QueryUpdateUserBalance, price, userID)
	if err != nil {
		slog.Error("Ошибка обновления баланса", "error", err)
		return fmt.Errorf("tx.Exec: %w", err)
	}

	_, err = tx.Exec(ctx, QueryUpdateInventory, userID, itemID)
	if err != nil {
		slog.Error("Ошибка обновления инвентаря", "error", err)
		return fmt.Errorf("tx.Exec: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при коммите транзакции", "error", err)
		return fmt.Errorf("tx.Commit: %w", err)
	}

	slog.Info("Покупка успешна", "userID", userID, "item", itemName, "price", price)
	return nil
}
