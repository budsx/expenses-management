package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/budsx/expenses-management/entity"
)

type ExpensesRepository struct {
	db *sql.DB
}

func NewExpensesRepository(db *sql.DB) *ExpensesRepository {
	return &ExpensesRepository{db: db}
}

func (r *ExpensesRepository) WriteExpense(ctx context.Context, expense *entity.Expense) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO expenses (user_id, amount_idr, description, receipt_url, status, submitted_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	now := time.Now()
	var id int64
	err = tx.QueryRowContext(
		ctx,
		query,
		expense.UserID,
		expense.AmountIDR,
		expense.Description,
		expense.ReceiptURL,
		expense.Status,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ExpensesRepository) ApprovalExpense(ctx context.Context, expenseApproval *entity.ExpenseApproval) (int64, error) {
	queryExpense := `
		UPDATE expenses SET status = $1 WHERE id = $2
	`

	queryApproval := `
		INSERT INTO approvals (expense_id, approver_id, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		queryExpense,
		expenseApproval.Status,
		expenseApproval.ExpenseID,
	)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(
		ctx,
		queryApproval,
		expenseApproval.ExpenseID,
		expenseApproval.ApproverID,
		expenseApproval.Status,
		expenseApproval.Notes,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return expenseApproval.ExpenseID, nil
}

func (r *ExpensesRepository) UpdateExpenseStatus(ctx context.Context, expenseID int64, status int32) error {
	query := `
		UPDATE expenses SET status = $1 WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, expenseID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ExpensesRepository) GetExpenseByID(ctx context.Context, expenseID int64) (*entity.Expense, error) {
	query := `
		SELECT id, user_id, amount_idr, description, receipt_url, status, auto_approved, submitted_at, processed_at FROM expenses WHERE id = $1
	`

	var expense entity.Expense
	err := r.db.QueryRowContext(ctx, query, expenseID).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.AmountIDR,
		&expense.Description,
		&expense.ReceiptURL,
		&expense.Status,
		&expense.AutoApproved,
		&expense.SubmittedAt,
		&sql.NullTime{},
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("expense", expense)
	return &expense, nil
}
