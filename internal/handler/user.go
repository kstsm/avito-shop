package handler

import (
	"encoding/json"
	"errors"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/repository"
	"net/http"
)

func (h Handler) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		http.Error(w, "не удалось получить userID", http.StatusUnauthorized)
		return
	}

	info, err := h.service.GetUserInfo(ctx, userID)
	if err != nil {
		slog.Error("Ошибка при получении информации", "error", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusOK, info)
}

func (h Handler) sendCoinsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	senderID, ok := ctx.Value("userID").(int)
	if !ok {
		slog.Warn("Неавторизованный запрос на отправку монет")
		writeErrorResponse(w, http.StatusUnauthorized, "Неавторизованный пользователь")
		return
	}

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Ошибка разбора JSON", "error", err.Error())
		writeErrorResponse(w, http.StatusBadRequest, "Некорректное тело запроса")
		return
	}
	defer r.Body.Close()

	if req.Amount <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Сумма должна быть больше 0")
		return
	}

	err := h.service.SendCoins(ctx, req.Amount, senderID, req.ToUser)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInsufficientFunds):
			slog.Warn("Недостаточно средств", "userID", senderID)
			writeErrorResponse(w, http.StatusBadRequest, "Недостаточно средств")

		case errors.Is(err, repository.ErrUserNotFound):
			slog.Warn("Получатель не найден", "recipient", req.ToUser)
			writeErrorResponse(w, http.StatusBadRequest, "Получатель не найден")

		case errors.Is(err, repository.ErrInvalidTransfer):
			slog.Warn("Нельзя перевести монеты самому себе", "recipient", req.ToUser)
			writeErrorResponse(w, http.StatusBadRequest, "Нельзя перевести монеты самому себе")

		default:
			slog.Error("Внутренняя ошибка сервера", "error", err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, "Успешный ответ.")
}
