package handler

import (
	"context"
	"encoding/json"
	"github.com/kstsm/avito-shop/api/rest/models"
	"net/http"
)

func (h Handler) authHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var req models.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"errors": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	token, err := h.service.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, `{"errors": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}
