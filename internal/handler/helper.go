package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kstsm/avito-shop/api/rest/models"
	"log"
	"net/http"
	"strings"
)

func parseAndValidateRequest(r *http.Request, req interface{}) *models.ErrorResponse {
	validate := validator.New()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil {
		return &models.ErrorResponse{Errors: "Некорректный JSON или неизвестные поля"}
	}

	if err := validate.Struct(req); err != nil {
		return formatValidationErrors(err)
	}

	return nil
}

func formatValidationErrors(err error) *models.ErrorResponse {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, formatErrorMessage(e))
		}
		return &models.ErrorResponse{Errors: strings.Join(errorMessages, ", ")}
	}
	return &models.ErrorResponse{Errors: "Ошибка валидации"}
}

func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Поле %s обязательно", e.Field())
	case "min":
		return fmt.Sprintf("Поле %s должно быть не короче %s символов", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("Поле %s должно быть не длиннее %s символов", e.Field(), e.Param())
	case "alphanum":
		return fmt.Sprintf("Поле %s должно содержать только буквы и цифры", e.Field())
	case "excludesrune":
		return fmt.Sprintf("Поле %s не должно содержать пробелы", e.Field())
	default:
		return fmt.Sprintf("Поле %s не прошло валидацию", e.Field())
	}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
	}
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Errors: message})
}
