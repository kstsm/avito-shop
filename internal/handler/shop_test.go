package handler

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/kstsm/avito-shop/internal/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuyItemHandler_Success(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	userID := 1
	item := "Cup"

	mockService.On("BuyItem", mock.Anything, userID, item).Return(nil)

	ctx := context.WithValue(context.Background(), "userID", userID)

	req := httptest.NewRequest("POST", "/buy/"+item, nil).WithContext(ctx)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"item"}, Values: []string{item}},
	}))

	rr := httptest.NewRecorder()
	handler.buyItemHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertCalled(t, "BuyItem", mock.Anything, userID, item)
}

func TestBuyItemHandler_InsufficientFunds(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	userID := 100
	item := "Phone"

	mockService.On("BuyItem", mock.Anything, userID, item).Return(apperrors.ErrInsufficientFunds)

	ctx := context.WithValue(context.Background(), "userID", userID)

	req := httptest.NewRequest("POST", "/buy/"+item, nil).WithContext(ctx)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"item"}, Values: []string{item}},
	}))

	rr := httptest.NewRecorder()
	handler.buyItemHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Недостаточно средств")
	mockService.AssertCalled(t, "BuyItem", mock.Anything, userID, item)
}

func TestBuyItemHandler_ItemNotFound(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	userID := 4
	item := "PC"

	mockService.On("BuyItem", mock.Anything, userID, item).Return(apperrors.ErrItemNotFound)

	ctx := context.WithValue(context.Background(), "userID", userID)

	req := httptest.NewRequest("POST", "/buy/"+item, nil).WithContext(ctx)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"item"}, Values: []string{item}},
	}))

	rr := httptest.NewRecorder()
	handler.buyItemHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Товар не найден")
	mockService.AssertCalled(t, "BuyItem", mock.Anything, userID, item)
}

func TestBuyItemHandler_InternalServerError(t *testing.T) {
	mockService := new(MockService)
	handler := Handler{service: mockService}

	userID := 5
	item := "Keyboard"

	mockService.On("BuyItem", mock.Anything, userID, item).Return(errors.New("db error"))

	ctx := context.WithValue(context.Background(), "userID", userID)

	req := httptest.NewRequest("POST", "/buy/"+item, nil).WithContext(ctx)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"item"}, Values: []string{item}},
	}))

	rr := httptest.NewRecorder()
	handler.buyItemHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Внутренняя ошибка сервера")
	mockService.AssertCalled(t, "BuyItem", mock.Anything, userID, item)
}
