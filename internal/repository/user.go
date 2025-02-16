package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/apperrors"
)

func (r Repository) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	var bytes []byte
	var result models.InfoResponse

	err := r.conn.QueryRow(ctx, QueryGetUserInfo, userID).Scan(&bytes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.InfoResponse{}, apperrors.ErrUserNotFound
		}
		slog.Error("Ошибка при запросе данных пользователя", "userID", userID, "error", err)
		return models.InfoResponse{}, fmt.Errorf("r.conn.QueryRow: %w", err)
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		slog.Error("Ошибка при декодировании JSON", "userID", userID, "error", err)
		return models.InfoResponse{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return result, nil
}

func (r Repository) SendCoins(ctx context.Context, senderID, amount int, username string) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		slog.Error("Ошибка при начале транзакции", "error", err)
		return fmt.Errorf("r.conn.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	balance, err := r.getUserBalance(ctx, tx, senderID)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrInsufficientFunds, err)
	}
	if balance < amount {
		return fmt.Errorf("%w: недостаточно средств", apperrors.ErrInsufficientFunds)
	}

	receiverID, err := r.getUserIDByUsername(ctx, tx, username)
	if err != nil {
		slog.Warn("Получатель не найден", username)
		return fmt.Errorf("%w", apperrors.ErrUserNotFound)
	}

	if senderID == receiverID {
		slog.Warn("Попытка перевода самому себе", "userID", senderID)
		return fmt.Errorf("%w: нельзя перевести монеты самому себе", apperrors.ErrInvalidTransfer)
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $2 WHERE id = $1", senderID, amount)
	if err != nil {
		slog.Error("Ошибка при обновлении баланса отправителя", "senderID", senderID, "error", err)
		return fmt.Errorf("Ошибка при обновлении баланса отправителя: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $2 WHERE id = $1", receiverID, amount)
	if err != nil {
		slog.Error("Ошибка при обновлении баланса получателя", "receiverID", receiverID, "error", err)
		return fmt.Errorf("Ошибка при обновлении баланса получателя: %w", err)
	}

	_, err = tx.Exec(ctx, QueryInsertTransaction, senderID, receiverID, amount)
	if err != nil {
		slog.Error("Ошибка при записи транзакции", "senderID", senderID, "receiverID", receiverID, "error", err)
		return fmt.Errorf("tx.Exec: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при коммите транзакции", "error", err)
		return fmt.Errorf("tx.Commit: %w", err)
	}

	slog.Info("Перевод выполнен", "from", senderID, "to", receiverID, "amount", amount)
	return nil
}
