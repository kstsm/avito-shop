package handler

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"net/http"
)

func (h Handler) buyItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	itemName := chi.URLParam(r, "item")

	err := h.service.BuyItem(h.ctx, userID, itemName)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInsufficientFunds):
			slog.Warn("Недостаточно средств", "userID", userID, "item", itemName)
			WriteErrorResponse(w, http.StatusBadRequest, "Недостаточно средств")
		case errors.Is(err, apperrors.ErrItemNotFound):
			slog.Warn("Товар не найден", "userID", userID, "item", itemName)
			WriteErrorResponse(w, http.StatusNotFound, "Товар не найден")
		default:
			slog.Error("Внутренняя ошибка сервера", "userID", userID, "item", itemName, "error", err)
			WriteErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	slog.Info("Покупка завершена", "userID", userID, "item", itemName)
	sendJSONResponse(w, http.StatusOK, nil)
}
