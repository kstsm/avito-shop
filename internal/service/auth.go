package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"github.com/kstsm/avito-shop/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if errors.Is(err, pgx.ErrNoRows) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Ошибка хеширования пароля", "username", username, "error", err.Error())
			return "", fmt.Errorf("не удалось хешировать пароль: %w", err)
		}

		userID, err := s.repo.CreateUser(ctx, username, string(hashedPassword))
		if err != nil {
			slog.Error("Ошибка создания пользователя", "username", username, "error", err.Error())
			return "", fmt.Errorf("%w: ошибка базы данных", apperrors.ErrUserCreation)
		}

		slog.Infof("Создан новый пользователь: '%s' с ID %d", username, userID)
		return auth.GenerateToken(userID)
	}

	if err != nil {
		slog.Error("Ошибка при получении пользователя", "username", username, "error", err.Error())
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		slog.Warn("Ошибка аутентификации: неверный пароль", "username", username)
		return "", apperrors.ErrInvalidCredentials
	}

	slog.Infof("Пользователь '%s' успешно аутентифицирован", username)
	return auth.GenerateToken(user.ID)
}
