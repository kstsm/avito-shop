package handler

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/kstsm/avito-shop/internal/auth"
	"github.com/kstsm/avito-shop/internal/middleware"
	"github.com/kstsm/avito-shop/internal/service"
	"net/http"
)

type HandlerI interface {
	NewRouter() http.Handler
	authHandler(w http.ResponseWriter, r *http.Request)
	getUserInfoHandler(w http.ResponseWriter, r *http.Request)
	sendCoinsHandler(w http.ResponseWriter, r *http.Request)
	buyItemHandler(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	ctx     context.Context
	service service.ServiceI
}

func NewHandler(ctx context.Context, svc service.ServiceI) HandlerI {
	return &Handler{
		ctx:     ctx,
		service: svc,
	}
}

func NewRouterForTests(ctx context.Context, svc service.ServiceI) http.Handler {
	router := NewHandler(ctx, svc)
	return router.NewRouter()
}

func (h Handler) NewRouter() http.Handler {
	router := chi.NewRouter()

	authMiddleware := middleware.AuthMiddleware(auth.ValidateToken)

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", h.authHandler)

		r.With(authMiddleware).Group(func(r chi.Router) {
			r.Get("/buy/{item}", h.buyItemHandler)
			r.Post("/sendCoin", h.sendCoinsHandler)
			r.Get("/info", h.getUserInfoHandler)

		})
	})

	return router
}
