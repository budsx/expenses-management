package repository

import (
	iface "github.com/budsx/expenses-management/repository/interface"
)

type Repository struct {
	PaymentProcessor   iface.PaymentProcessor
	UserRepository     iface.UserRepository
	ExpensesRepository iface.ExpensesRepository
}

func NewRepository(paymentProcessor iface.PaymentProcessor, userRepository iface.UserRepository, expensesRepository iface.ExpensesRepository) *Repository {
	return &Repository{
		PaymentProcessor:   paymentProcessor,
		UserRepository:     userRepository,
		ExpensesRepository: expensesRepository,
	}
}
