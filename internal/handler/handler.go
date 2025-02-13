package handler

import (
	"github.com/go-chi/chi"
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
	service service.ServiceI
}

func NewHandler(svc service.ServiceI) HandlerI {
	return &Handler{
		service: svc,
	}
}

func (h Handler) NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", h.authHandler)

		r.With(middleware.AuthMiddleware()).Group(func(r chi.Router) {
			r.Get("/info", h.getUserInfoHandler)
			r.Get("/buy/{item}", h.buyItemHandler)
			r.Post("/sendCoin", h.sendCoinsHandler)
		})
	})

	return router
}
