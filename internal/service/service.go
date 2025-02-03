package service

import (
	"context"
	"fmt"
	"project1/internal/models"
	"project1/internal/repository"
)

type Operations interface {
	Deposit(ctx context.Context, id int, amount float64) error
	Transfer(ctx context.Context, fromId int, toId int, amount float64) error
	GetLastTransactions(ctx context.Context, userId int) ([]models.Transaction, error)
}

type FinancialOperator struct {
	rep repository.Repository
}

func NewFinancialOperator(rep repository.Repository) *FinancialOperator {
	return &FinancialOperator{
		rep: rep,
	}
}

func (fo *FinancialOperator) Deposit(ctx context.Context, id int, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount can not be less than 0")
	}
	var amountInt int
	amountInt = int(amount * 100)
	err := fo.rep.Deposit(ctx, id, amountInt)
	if err != nil {
		return err
	}
	return nil
}

func (fo *FinancialOperator) Transfer(ctx context.Context, fromId int, toId int, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount can not be less than 0")
	}
	if fromId == toId {
		return fmt.Errorf("same sander and reciever")
	}

	var amountInt int
	amountInt = int(amount * 100)
	err := fo.rep.Transfer(ctx, fromId, toId, amountInt)
	if err != nil {
		return err
	}
	return nil
}

func (fo *FinancialOperator) GetLastTransactions(ctx context.Context, userId int) ([]models.Transaction, error) {
	res, err := fo.rep.GetLastTransactions(ctx, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
