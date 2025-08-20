package postgres

import "database/sql"

type ExpensesRepository struct {
	db *sql.DB
}

func NewExpensesRepository(db *sql.DB) *ExpensesRepository {
	return &ExpensesRepository{db: db}
}
