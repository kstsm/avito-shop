package handler

import (
	"context"
	"fmt"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Authenticate(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	fmt.Println("Mock GetUserInfo called") // Проверка, вызывается ли метод
	args := m.Called(ctx, userID)

	var response models.InfoResponse
	if args.Get(0) != nil {
		response = args.Get(0).(models.InfoResponse)
	}

	return response, args.Error(1)
}

func (m *MockService) BuyItem(ctx context.Context, userID int, item string) error {
	args := m.Called(ctx, userID, item)
	return args.Error(0)
}

func (m *MockService) SendCoins(ctx context.Context, amount int, senderID int, username string) error {
	args := m.Called(ctx, amount, senderID, username)
	return args.Error(0)
}
