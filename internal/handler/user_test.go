package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserInfoHandler_Success(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{
		service: mockService,
	}

	userID := 6
	expectedResult := models.InfoResponse{
		Coins:     1000,
		Inventory: []models.Item{{Type: "cup", Quantity: 1}},
		CoinHistory: models.CoinHistory{
			Received: []models.ReceivedTransaction{
				{Amount: 1, FromUser: "User1"},
			},
			Sent: []models.SentTransaction{
				{Amount: 1, ToUser: "User2"},
			},
		},
	}

	mockService.On("GetUserInfo", mock.Anything, userID).Return(expectedResult, nil)
	req, err := http.NewRequest("GET", "/api/info", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context.WithValue(req.Context(), "userID", userID))

	rr := httptest.NewRecorder()
	handler.getUserInfoHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	data := map[string]interface{}{
		"coins": 1000,
		"inventory": []map[string]interface{}{
			{"type": "cup", "quantity": 1},
		},
		"coinHistory": map[string]interface{}{
			"received": []map[string]interface{}{
				{"amount": 1, "fromUser": "User1"},
			},
			"sent": []map[string]interface{}{
				{"amount": 1, "toUser": "User2"},
			},
		},
	}

	expectedResponse, err := json.Marshal(data)
	if err != nil {
		slog.Info("Ошибка")
		return
	}

	assert.JSONEq(t, string(expectedResponse), rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestGetUserInfoHandler_InternalServerError(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{
		service: mockService,
	}

	userID := 1

	mockService.On("GetUserInfo", mock.Anything, userID).Return(models.InfoResponse{}, errors.New("some unknown error"))

	req, err := http.NewRequest("GET", "/api/info", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context.WithValue(req.Context(), "userID", userID))

	rr := httptest.NewRecorder()
	handler.getUserInfoHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponse := `{"errors":"Внутренняя ошибка сервера"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestSendCoinsHandler_Success(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	senderID := 1
	reqBody := models.SendCoinRequest{ToUser: "user2", Amount: 100}
	mockService.On("SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/sendCoins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", senderID))

	rr := httptest.NewRecorder()
	handler.sendCoinsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertCalled(t, "SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser)
}

func TestSendCoinsHandler_ValidationError(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	senderID := 1
	invalidBody := `{"toUser": "Banner", "amount": -10}` // Ошибки: toUser пустой, amount < 0
	req := httptest.NewRequest("GET", "/api/sendCoin", bytes.NewReader([]byte(invalidBody)))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", senderID))

	rr := httptest.NewRecorder()
	handler.sendCoinsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Ошибка валидации запроса")
	mockService.AssertNotCalled(t, "SendCoins")
}

func TestSendCoinsHandler_InsufficientFunds(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	senderID := 1
	reqBody := models.SendCoinRequest{ToUser: "user2", Amount: 5000}
	mockService.On("SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser).Return(apperrors.ErrInsufficientFunds)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/sendCoins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", senderID))

	rr := httptest.NewRecorder()
	handler.sendCoinsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Недостаточно средств")
	mockService.AssertCalled(t, "SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser)
}

func TestSendCoinsHandler_UserNotFound(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	senderID := 1
	reqBody := models.SendCoinRequest{ToUser: "unknown_user", Amount: 50}
	mockService.On("SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser).Return(apperrors.ErrUserNotFound)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/sendCoins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", senderID))

	rr := httptest.NewRecorder()
	handler.sendCoinsHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Получатель не найден")
	mockService.AssertCalled(t, "SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser)
}

func TestSendCoinsHandler_InvalidTransfer(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	senderID := 1
	reqBody := models.SendCoinRequest{ToUser: "user1", Amount: 50} // Отправка самому себе
	mockService.On("SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser).Return(apperrors.ErrInvalidTransfer)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/sendCoins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", senderID))

	rr := httptest.NewRecorder()
	handler.sendCoinsHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Нельзя перевести монеты самому себе")
	mockService.AssertCalled(t, "SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser)
}

func TestSendCoinsHandler_InternalServerError(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	senderID := 1
	reqBody := models.SendCoinRequest{ToUser: "user2", Amount: 100}
	mockService.On("SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser).Return(errors.New("database error"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/sendCoins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", senderID))

	rr := httptest.NewRecorder()
	handler.sendCoinsHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Внутренняя ошибка сервера")
	mockService.AssertCalled(t, "SendCoins", mock.Anything, reqBody.Amount, senderID, reqBody.ToUser)
}
