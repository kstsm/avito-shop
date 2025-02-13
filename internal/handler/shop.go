package handler

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/internal/repository"
	"net/http"
)

func (h Handler) buyItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "Неавторизованный пользователь")
		return
	}

	itemName := chi.URLParam(r, "item")
	if itemName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректное имя предмета")
		return
	}

	err := h.service.BuyItem(ctx, userID, itemName)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInsufficientFunds):
			slog.Warn("Недостаточно средств", "userID", userID)
			writeErrorResponse(w, http.StatusBadRequest, "Недостаточно средств")

		case errors.Is(err, repository.ErrUserNotFound):
			slog.Warn("Пользователь не найден", "userID", userID)
			writeErrorResponse(w, http.StatusBadRequest, "Пользователь не найден")

		case errors.Is(err, repository.ErrItemNotFound):
			slog.Warn("Товар не найден", "itemName", itemName)
			writeErrorResponse(w, http.StatusBadRequest, "Товар не найден")

		default:
			slog.Error("Внутренняя ошибка сервера", "error", err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, "Успешный ответ.")
}
