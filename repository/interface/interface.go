package _interface

import (
	"context"

	"github.com/budsx/expenses-management/model"
)

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, payment *model.PaymentProcessorModel) error
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type ExpensesRepository interface {
}
