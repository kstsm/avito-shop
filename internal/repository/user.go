package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/api/rest/models"
)

func (r Repository) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	var bytes []byte
	var infoUser models.InfoResponse

	err := r.conn.QueryRow(ctx, QueryGetUserInfo, userID).Scan(&bytes)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("ошибка при запросе данных: %w", err)
	}

	err = json.Unmarshal(bytes, &infoUser)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	return infoUser, nil
}

func (r Repository) TransferCoins(ctx context.Context, senderID, amount int, username string) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		slog.Error("Ошибка при начале транзакции", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}
	defer tx.Rollback(ctx)

	balance, err := r.getUserBalance(ctx, tx, senderID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInsufficientFunds, err)
	}
	if balance < amount {
		return fmt.Errorf("%w: недостаточно средств", ErrInsufficientFunds)
	}

	receiverID, err := r.getUserIDByUsername(ctx, tx, username)
	if err != nil {
		return ErrUserNotFound
	}

	if senderID == receiverID {
		slog.Warn("Попытка перевода самому себе", "userID", senderID)
		return fmt.Errorf("%w: нельзя перевести монеты самому себе", ErrInvalidTransfer)
	}

	_, err = tx.Exec(ctx, QueryTransferCoins, senderID, amount, receiverID)
	if err != nil {
		slog.Error("Ошибка при обновлении баланса", "senderID", senderID, "receiverID", receiverID, "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	_, err = tx.Exec(ctx, QueryInsertTransaction, senderID, receiverID, amount)
	if err != nil {
		slog.Error("Ошибка при записи транзакции", "senderID", senderID, "receiverID", receiverID, "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при коммите транзакции", "error", err)
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	slog.Info("Перевод выполнен", "from", senderID, "to", receiverID, "amount", amount)
	return nil
}
