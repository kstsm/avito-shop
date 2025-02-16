package service

import (
	"context"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) GenerateToken(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) BuyItem(ctx context.Context, userID int, itemName string) error {
	args := m.Called(ctx, userID, itemName)
	return args.Error(0)
}

func (m *MockRepository) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.InfoResponse), args.Error(1)
}

func (m *MockRepository) SendCoins(ctx context.Context, senderID int, amount int, toUser string) error {
	args := m.Called(ctx, senderID, amount, toUser)
	return args.Error(0)
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, username string, hashedPassword string) (int, error) {
	args := m.Called(ctx, username, hashedPassword)
	return args.Int(0), args.Error(1)
}
