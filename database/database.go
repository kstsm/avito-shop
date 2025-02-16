package database

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/avito-shop/config"
	"log"
)

type Repository struct {
	dg *pgxpool.Pool
}

var cfg = config.Config

func InitPostgres(ctx context.Context) *pgxpool.Pool {
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

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	fmt.Println("Успешное подключение к базе данных")

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Ошибка проверки соединения с БД: %v", err)
	}
	fmt.Println("База данных доступна")

	return pool
}

func InitTestPostgres(ctx context.Context) (*pgxpool.Pool, error) {
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

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	slog.Info("Успешное подключение к базе данных")
	return pool, nil
}
