package handler

import (
	"github.com/go-chi/chi"
	"github.com/kstsm/avito-shop/internal/service"
	"net/http"
)

type HandlerI interface {
	getApiInfoHandler(w http.ResponseWriter, r *http.Request)
	NewRouter() http.Handler
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
		
		r.Post("/info", h.getApiInfoHandler)
	})

	/*router.Route("/house", func(r chi.Router) {
		r.With(h.AuthMW.JWTAuth).Post("/create", h.CreateHouseHandler)
	})*/

	return router
}
