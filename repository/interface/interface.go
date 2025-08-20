package _interface

import (
	"context"

	"github.com/budsx/expenses-management/model"
)

type PaymentProcessor interface {
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
}

type ExpensesRepository interface {
}
