package database

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/config"
)

type Repository struct {
	dg *pgx.Conn
}

func InitPostgres() *pgx.Conn {
	cfg := config.Config.Postgres

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		slog.Fatal("Ошибка подключения к базе данных: %v", err)
	}

	return conn
}
