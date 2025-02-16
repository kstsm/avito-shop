package tests

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/avito-shop/database"
	"github.com/kstsm/avito-shop/internal/handler"
	"github.com/kstsm/avito-shop/internal/repository"
	"github.com/kstsm/avito-shop/internal/service"
	"net/http/httptest"
	"testing"
)

func SetupTestServer(t *testing.T) (*httptest.Server, context.Context, *pgx.Conn) {
	ctx := context.Background()
	conn, err := database.InitTestPostgres(ctx)
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		t.Fatalf("Ошибка при проверке подключения к базе данных: %v", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(ctx); err != nil {
			t.Errorf("Ошибка закрытия соединения с базой данных: %v", err)
		}
	})

	repo := repository.NewRepository(conn)
	svc := service.NewService(repo)
	router := handler.NewRouterForTests(ctx, svc)
	ts := httptest.NewServer(router)

	return ts, ctx, conn
}
