package cmd

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/config"
	"github.com/kstsm/avito-shop/database"
	"github.com/kstsm/avito-shop/internal/handler"
	"github.com/kstsm/avito-shop/internal/repository"
	"github.com/kstsm/avito-shop/internal/service"
	"net/http"
)

func Run() {
	cfg := config.Config

	db := database.InitPostgres()

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	router := handler.NewHandler(svc)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router.NewRouter(),
	}

	slog.Info("Сервер запущен", "host", cfg.Server.Host, "port", cfg.Server.Port)
	err := srv.ListenAndServe()
	if err != nil {
		slog.Fatal("Ошибка при запуске сервера", "error", err)
	}
}
