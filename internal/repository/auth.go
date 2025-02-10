package repository

import (
	"context"
	"errors"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/api/rest/models"
)

func (r Repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	err := r.conn.QueryRow(ctx, GetUserByUsername, username).Scan(&user.ID, &user.Username, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Warnf("Пользователь '%s' не найден", username)
		return models.User{}, pgx.ErrNoRows
	}
	if err != nil {
		slog.Errorf("Ошибка при запросе пользователя '%s': %v", username, err)
		return models.User{}, err
	}

	slog.Infof("Пользователь '%s' успешно найден", username)
	return user, nil
}

func (r Repository) CreateUser(ctx context.Context, username, hashedPassword string) error {
	_, err := r.conn.Exec(ctx, CreateUser, username, hashedPassword)
	if err != nil {
		slog.Errorf("Ошибка при добавлении пользователя '%s': %v", username, err)
		return err
	}

	slog.Infof("Пользователь '%s' успешно добавлен в базу данных", username)
	return nil
}
