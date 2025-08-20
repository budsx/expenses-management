package service

import (
	repo "github.com/budsx/expenses-management/repository"
	"github.com/sirupsen/logrus"
)

type ExpensesManagementService struct {
	repository *repo.Repository
	logger     *logrus.Logger
}

func NewExpensesManagementService(repository *repo.Repository, logger *logrus.Logger) *ExpensesManagementService {
	return &ExpensesManagementService{repository: repository, logger: logger}
}
