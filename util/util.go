package util

import (
	"context"
	"fmt"

	"github.com/budsx/expenses-management/model"
)

func GetUserInfoFromContext(ctx context.Context) (model.User, error) {
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		fmt.Println("user_id not found")
		return model.User{}, fmt.Errorf("user_id not found")
	}

	email, ok := ctx.Value("user_email").(string)
	if !ok {
		return model.User{}, fmt.Errorf("email not found")
	}

	role, ok := ctx.Value("user_role").(int)
	if !ok {
		return model.User{}, fmt.Errorf("role not found")
	}

	return model.User{
		ID:    userID,
		Email: email,
		Role:  role,
	}, nil
}
