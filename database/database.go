package database

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/config"
	"os"
)

type Repository struct {
	dg *pgx.Conn
}

var cfg = config.Config

func InitPostgres(ctx context.Context) *pgx.Conn {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
	)

	slog.Info(
		"Подключение к базе данных", "host", cfg.Postgres.Host,
		"port", cfg.Postgres.Port, "db", cfg.Postgres.DBName,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		slog.Fatal("Ошибка подключения к базе данных: %v", err)
		os.Exit(1)
	}

	slog.Info("Успешное подключение к базе данных")
	return conn
}

func InitTestPostgres(ctx context.Context) (*pgx.Conn, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		"test_user",
		"test_pass",
		"localhost",
		"5433",
		"test_db",
	)

	slog.Info(
		"Подключение к базе данных", "host", cfg.Postgres.Host,
		"port", cfg.Postgres.Port, "db", cfg.Postgres.DBName,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		slog.Fatal("Ошибка подключения к базе данных", "error", err)
		return nil, err
	}

	slog.Info("Успешное подключение к базе данных")
	return conn, nil
}
