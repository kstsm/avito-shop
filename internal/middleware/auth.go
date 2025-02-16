package middleware

import (
	"context"
	"encoding/json"
	"github.com/gookit/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(validateToken func(string) (int, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ExtractToken(r)
			if token == "" {
				sendJSONError(w, http.StatusUnauthorized, "Отсутствует токен авторизации")
				return
			}
			userID, err := validateToken(token)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "Неверный или просроченный токен")
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		slog.Warn("Заголовок Authorization отсутствует", "path", r.URL.Path)
		return ""
	}
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		slog.Warn("Некорректный формат токена", "path", r.URL.Path, "token", bearerToken)
		return ""
	}

	return parts[1]
}

func sendJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := struct {
		Errors string `json:"errors"`
	}{
		Errors: message,
	}

	json.NewEncoder(w).Encode(errorResponse)
}
