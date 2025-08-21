package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/budsx/expenses-management/entity"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserWithPassword(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = $1`

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
