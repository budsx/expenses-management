package repository

import (
	iface "github.com/budsx/expenses-management/repository/interface"
)

type Repository struct {
	PaymentProcessor   iface.PaymentProcessor
	UserRepository     iface.UserRepository
	ExpensesRepository iface.ExpensesRepository
	RabbitMQClient     iface.RabbitMQClient
}

func NewRepository(paymentProcessor iface.PaymentProcessor, userRepository iface.UserRepository, expensesRepository iface.ExpensesRepository, rabbitmqClient iface.RabbitMQClient) *Repository {
	return &Repository{
		PaymentProcessor:   paymentProcessor,
		UserRepository:     userRepository,
		ExpensesRepository: expensesRepository,
		RabbitMQClient:     rabbitmqClient,
	}
}
