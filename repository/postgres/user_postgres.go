package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/budsx/expenses-management/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	fmt.Println("Testing GetUser Repository")
	return &model.User{
		ID:        id,
		Username:  "test",
		Email:     "test@test.com",
		CreatedAt: time.Now(),
	}, nil
}
