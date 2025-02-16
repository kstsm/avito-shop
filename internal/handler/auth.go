package handler

import (
	"errors"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"net/http"
)

func (h Handler) authHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		slog.Warn("Ошибка валидации запроса", "error", err)
		WriteErrorResponse(w, http.StatusBadRequest, err.Errors)
		return
	}

	token, err := h.service.Authenticate(h.ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			slog.Warn("Ошибка аутентификации", "username", req.Username, "error", err.Error())
			WriteErrorResponse(w, http.StatusUnauthorized, "Неверное имя пользователя или пароль")
		default:
			slog.Error("Внутренняя ошибка сервера", "error", err.Error())
			WriteErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, models.AuthResponse{Token: token})
}
