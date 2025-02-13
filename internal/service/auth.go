package service

import (
	"context"
	"errors"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)

	if errors.Is(err, pgx.ErrNoRows) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Errorf("Ошибка хеширования пароля для '%s': %v", username, err)
			return "", errors.New("failed to hash password")
		}

		userID, err := s.repo.CreateUser(ctx, username, string(hashedPassword))
		if err != nil {
			slog.Errorf("Ошибка создания пользователя '%s': %v", username, err)
			return "", err
		}

		slog.Infof("Создан новый пользователь: '%s' с ID %d", username, userID)
		return auth.GenerateToken(userID)
	}

	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		slog.Warnf("Ошибка аутентификации для '%s': неверный пароль", username)
		return "", errors.New("invalid credentials")
	}

	slog.Infof("Пользователь '%s' успешно аутентифицирован", username)
	return auth.GenerateToken(user.ID)
}
