package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuyItem_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := Service{repo: mockRepo}

	mockRepo.On("BuyItem", mock.Anything, 1, "T-shirt").Return(nil)

	err := service.BuyItem(context.Background(), 1, "T-shirt")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBuyItem_InsufficientFunds(t *testing.T) {
	mockRepo := new(MockRepository)
	service := Service{repo: mockRepo}

	mockRepo.On("BuyItem", mock.Anything, 1, "T-shirt").Return(errors.New("insufficient funds"))

	err := service.BuyItem(context.Background(), 1, "T-shirt")

	assert.Error(t, err)
	assert.EqualError(t, err, "insufficient funds")

	mockRepo.AssertExpectations(t)
}
