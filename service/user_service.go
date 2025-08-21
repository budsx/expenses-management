package service

import (
	"context"
	"fmt"
	"time"

	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/util"
)

func (s *ExpensesManagementService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repo.UserRepository.GetUser(ctx, id)
}

func (s *ExpensesManagementService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.UserRepository.GetUserByEmail(ctx, email)
}

func (s *ExpensesManagementService) AuthenticateUser(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	user, err := s.repo.UserRepository.GetUserWithPassword(ctx, email)
	if err != nil {
		s.logger.WithError(err).Error("invalid password verification")
		return nil, fmt.Errorf("invalid credentials")
	}

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		s.logger.WithError(err).Error("invalid password verification")
		return nil, fmt.Errorf("invalid credentials")
	}

	token, expiresAt, err := util.GenerateJWT(user.ID, user.Email, user.Role, time.Now().Add(24*time.Hour))
	if err != nil {
		s.logger.WithError(err).Error("failed to generate token")
		return nil, fmt.Errorf("failed to generate token")
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
