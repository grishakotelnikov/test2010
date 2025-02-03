package service_test

import (
	"context"
	"errors"
	"project1/internal/models"
	"project1/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Deposit(ctx context.Context, id int, amount int) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *MockRepository) Transfer(ctx context.Context, fromId int, toId int, amount int) error {
	args := m.Called(ctx, fromId, toId, amount)
	return args.Error(0)
}

func (m *MockRepository) GetLastTransactions(ctx context.Context, userId int) ([]models.Transaction, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func TestDeposit(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewFinancialOperator(mockRepo)
	ctx := context.Background()

	mockRepo.On("Deposit", ctx, 1, 1000).Return(nil)
	err := service.Deposit(ctx, 1, 10.0)
	assert.NoError(t, err)

	err = service.Deposit(ctx, 1, -10.0)
	assert.Error(t, err)
	assert.Equal(t, "amount can not be less than 0", err.Error())

	mockRepo.On("Deposit", ctx, 2, 500).Return(errors.New("db error"))
	err = service.Deposit(ctx, 2, 5.0)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestTransfer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewFinancialOperator(mockRepo)
	ctx := context.Background()

	mockRepo.On("Transfer", ctx, 1, 2, 1000).Return(nil)
	err := service.Transfer(ctx, 1, 2, 10.0)
	assert.NoError(t, err)

	err = service.Transfer(ctx, 1, 2, -10.0)
	assert.Error(t, err)
	assert.Equal(t, "amount can not be less than 0", err.Error())

	err = service.Transfer(ctx, 1, 1, 10.0)
	assert.Error(t, err)
	assert.Equal(t, "same sander and reciever", err.Error())

	mockRepo.On("Transfer", ctx, 1, 3, 500).Return(errors.New("db error"))
	err = service.Transfer(ctx, 1, 3, 5.0)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestGetLastTransactions(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewFinancialOperator(mockRepo)
	ctx := context.Background()

	transactions := []models.Transaction{
		{ID: 1, UserId: 1, Amount: 1000},
		{ID: 2, UserId: 1, Amount: -500},
	}

	mockRepo.On("GetLastTransactions", ctx, 1).Return(transactions, nil)
	res, err := service.GetLastTransactions(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, transactions, res)

	mockRepo.On("GetLastTransactions", ctx, 2).Return([]models.Transaction{}, errors.New("db error"))
	res, err = service.GetLastTransactions(ctx, 2)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "db error", err.Error())
}
