package handler

import (
	"encoding/json"
	"github.com/kstsm/avito-shop/api/rest/models"
	"log"
	"net/http"
)

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Errors: message})
}
