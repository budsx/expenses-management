package service

import (
	repo "github.com/budsx/expenses-management/repository"
	"github.com/sirupsen/logrus"
)

type ExpensesManagementService struct {
	repo   *repo.Repository
	logger *logrus.Logger
}

func NewExpensesManagementService(repo *repo.Repository, logger *logrus.Logger) *ExpensesManagementService {
	return &ExpensesManagementService{repo: repo, logger: logger}
}
