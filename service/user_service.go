package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/repository/postgres"
	"github.com/budsx/expenses-management/util"
)

func (s *ExpensesManagementService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repository.UserRepository.GetUser(ctx, id)
}

func (s *ExpensesManagementService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repository.UserRepository.GetUserByEmail(ctx, email)
}

func (s *ExpensesManagementService) AuthenticateUser(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	userRepo, ok := s.repository.UserRepository.(*postgres.UserRepository)
	if !ok {
		return nil, fmt.Errorf("invalid repository type")
	}

	user, err := userRepo.GetUserWithPassword(ctx, email)
	if err != nil {
		s.logger.Error("error getting user with password", "error", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		s.logger.Error("invalid password verification")
		return nil, fmt.Errorf("invalid credentials")
	}

	userIDStr := strconv.FormatInt(user.ID, 10)
	token, expiresAt, err := util.GenerateJWT(userIDStr, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
