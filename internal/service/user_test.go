package service

import (
	"context"
	"errors"
	"github.com/kstsm/avito-shop/api/rest/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserInfo(t *testing.T) {
	mockRepo := new(MockRepository)
	service := Service{repo: mockRepo}

	expectedResponse := models.InfoResponse{
		Coins: 1000,
		Inventory: []models.Item{
			{Type: "cup", Quantity: 2},
		},
		CoinHistory: models.CoinHistory{
			Received: []models.ReceivedTransaction{
				{FromUser: "Alice", Amount: 200},
			},
			Sent: []models.SentTransaction{
				{ToUser: "Bob", Amount: 100},
			},
		},
	}

	mockRepo.On("GetUserInfo", mock.Anything, 1).Return(expectedResponse, nil)

	response, err := service.GetUserInfo(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)

	mockRepo.AssertExpectations(t)
}

func TestSendCoins_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := Service{repo: mockRepo}

	mockRepo.On("SendCoins", mock.Anything, 1, 100, "Bob").Return(nil)

	err := service.SendCoins(context.Background(), 100, 1, "Bob")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSendCoins_InsufficientFunds(t *testing.T) {
	mockRepo := new(MockRepository)
	service := Service{repo: mockRepo}

	mockRepo.On("SendCoins", mock.Anything, 1, 5000, "Bob").Return(errors.New("insufficient funds"))

	err := service.SendCoins(context.Background(), 5000, 1, "Bob")

	assert.Error(t, err)
	assert.EqualError(t, err, "insufficient funds")

	mockRepo.AssertExpectations(t)
}
