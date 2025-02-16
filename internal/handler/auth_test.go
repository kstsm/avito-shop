package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Success(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	reqBody := models.AuthRequest{Username: "testuser", Password: "validPass123"}
	reqJSON, _ := json.Marshal(reqBody)

	expectedToken := "valid_token"
	mockService.On("Authenticate", mock.Anything, reqBody.Username, reqBody.Password).Return(expectedToken, nil)

	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.authHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := `{"token":"valid_token"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	mockService.AssertCalled(t, "Authenticate", mock.Anything, reqBody.Username, reqBody.Password)
}

func TestAuthHandler_InvalidCredentials(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	reqBody := models.AuthRequest{Username: "testuser", Password: "wrongPass"}
	reqJSON, _ := json.Marshal(reqBody)

	mockService.On("Authenticate", mock.Anything, reqBody.Username, reqBody.Password).Return("", apperrors.ErrInvalidCredentials)

	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.authHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponse := `{"errors":"Неверное имя пользователя или пароль"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	mockService.AssertCalled(t, "Authenticate", mock.Anything, reqBody.Username, reqBody.Password)
}

func TestAuthHandler_InvalidRequest(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	reqJSON := `{"username":"testuser"}`

	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer([]byte(reqJSON)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.authHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertNotCalled(t, "Authenticate")
}

func TestAuthHandler_InternalServerError(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	reqBody := models.AuthRequest{Username: "testuser", Password: "validPass123"}
	reqJSON, _ := json.Marshal(reqBody)

	mockService.On("Authenticate", mock.Anything, reqBody.Username, reqBody.Password).Return("", errors.New("db connection failed"))

	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.authHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponse := `{"errors":"Внутренняя ошибка сервера"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	mockService.AssertCalled(t, "Authenticate", mock.Anything, reqBody.Username, reqBody.Password)
}
