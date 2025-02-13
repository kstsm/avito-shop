package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
)

var (
	ErrInsufficientFunds = errors.New("недостаточно средств")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrItemNotFound      = errors.New("товар не найден")
	ErrTransactionFailed = errors.New("ошибка при выполнении транзакции")
	ErrInvalidTransfer   = errors.New("нельзя перевести монеты самому себе")
)

func (r Repository) BuyItem(ctx context.Context, userID int, itemName string) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		slog.Error("Ошибка при начале транзакции", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}
	defer tx.Rollback(ctx)

	var itemID, price int
	err = tx.QueryRow(ctx, QueryGetItem, itemName).Scan(&itemID, &price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Товар не найден", "item", itemName)
			return fmt.Errorf("%w: %s", ErrItemNotFound, itemName)
		}
		slog.Error("Ошибка при получении товара", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	balance, err := r.getUserBalance(ctx, tx, userID)
	if err != nil {
		slog.Error("Ошибка при получении баланса", "userID", userID, "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	if balance < price {
		slog.Warn("Недостаточно средств", "userID", userID, "balance", balance, "price", price)
		return ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, QueryUpdateUserBalance, price, userID)
	if err != nil {
		slog.Error("Ошибка обновления баланса", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	_, err = tx.Exec(ctx, QueryUpdateInventory, userID, itemID)
	if err != nil {
		slog.Error("Ошибка обновления инвентаря", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при коммите транзакции", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	slog.Info("Покупка успешна", "userID", userID, "item", itemName, "price", price)
	return nil
}
