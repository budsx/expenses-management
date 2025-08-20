package service

import (
	"context"

	"github.com/budsx/expenses-management/model"
)

func (s *ExpensesManagementService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repository.UserRepository.GetUser(ctx, id)
}
	