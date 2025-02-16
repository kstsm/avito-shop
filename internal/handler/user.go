package handler

import (
	"errors"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"net/http"
)

func (h Handler) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	result, err := h.service.GetUserInfo(h.ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			slog.Error("Пользователь не найден", "error", err)
			WriteErrorResponse(w, http.StatusBadRequest, "Пользователь не найден")
		default:
			slog.Error("Неизвестная ошибка", "error", err)
			WriteErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, result)
}

func (h Handler) sendCoinsHandler(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value("userID").(int)

	var req models.SendCoinRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		slog.Warn("Ошибка валидации запроса", "error", err)
		WriteErrorResponse(w, http.StatusBadRequest, "Ошибка валидации запроса")
		return
	}

	err := h.service.SendCoins(h.ctx, req.Amount, senderID, req.ToUser)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInsufficientFunds):
			slog.Warn("Недостаточно средств", "userID", senderID)
			WriteErrorResponse(w, http.StatusBadRequest, "Недостаточно средств")

		case errors.Is(err, apperrors.ErrInvalidTransfer):
			slog.Warn("Нельзя перевести монеты самому себе", "recipient", req.ToUser)
			WriteErrorResponse(w, http.StatusBadRequest, "Нельзя перевести монеты самому себе")

		case errors.Is(err, apperrors.ErrUserNotFound):
			slog.Warn("Получатель не найден", "recipient", req.ToUser)
			WriteErrorResponse(w, http.StatusNotFound, "Получатель не найден")

		default:
			slog.Error("Внутренняя ошибка сервера", "error", err.Error())
			WriteErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, nil)
}
