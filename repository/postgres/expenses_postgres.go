package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/budsx/expenses-management/entity"
)

type expensesRepository struct {
	db *sql.DB
}

func NewExpensesRepository(db *sql.DB) *expensesRepository {
	return &expensesRepository{db: db}
}

func (r *expensesRepository) PingContext(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *expensesRepository) WriteExpense(ctx context.Context, expense *entity.Expense) (int64, error) {
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

func (r *expensesRepository) ApprovalExpense(ctx context.Context, expenseApproval *entity.ExpenseApproval) error {
	queryExpense := `
		UPDATE expenses SET status = $1 WHERE id = $2
	`

	queryApproval := `
		INSERT INTO approvals (expense_id, approver_id, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		queryExpense,
		expenseApproval.Status,
		expenseApproval.ExpenseID,
	)
	if err != nil {
		return err
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
		return err
	}

	return tx.Commit()
}

func (r *expensesRepository) UpdateExpenseStatus(ctx context.Context, expenseID int64, status int32) error {
	query := `
		UPDATE expenses SET status = $1 WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, expenseID)
	if err != nil {
		return err
	}

	return nil
}

func (r *expensesRepository) GetExpenseByID(ctx context.Context, expenseID int64) (*entity.Expense, error) {
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

	return &expense, nil
}

func (r *expensesRepository) GetExpensesWithPagination(ctx context.Context, query *entity.ExpenseListQuery) ([]*entity.Expense, int64, error) {
	rows, err := r.db.QueryContext(ctx, buildDataQuery(query))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	expenses := make([]*entity.Expense, 0)

	sqlNullTime := sql.NullTime{}
	for rows.Next() {
		var expense entity.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.AmountIDR,
			&expense.Description,
			&expense.ReceiptURL,
			&expense.Status,
			&expense.AutoApproved,
			&expense.SubmittedAt,
			&sqlNullTime,
		)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, &expense)
	}

	totalCount := int64(0)
	err = r.db.QueryRowContext(ctx, buildQueryCount(query)).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return expenses, totalCount, nil
}

func buildDataQuery(query *entity.ExpenseListQuery) string {
	queryString := "SELECT id, user_id, amount_idr, description, receipt_url, status, auto_approved, submitted_at, processed_at FROM expenses"

	var conditions []string
	if query.UserID != 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = %d", query.UserID))
	}

	if query.Status != 0 {
		conditions = append(conditions, fmt.Sprintf("status = %d", query.Status))
	}

	if len(conditions) > 0 {
		queryString += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			queryString += " AND " + conditions[i]
		}
	}

	queryString += " ORDER BY id DESC"
	offset := (query.Page - 1) * query.Limit
	queryString += fmt.Sprintf(" LIMIT %d OFFSET %d", query.Limit, offset)

	return queryString
}

func buildQueryCount(query *entity.ExpenseListQuery) string {
	queryString := "SELECT COUNT(*) FROM expenses"

	var conditions []string
	if query.UserID != 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = %d", query.UserID))
	}

	if query.Status != 0 {
		conditions = append(conditions, fmt.Sprintf("status = %d", query.Status))
	}

	if len(conditions) > 0 {
		queryString += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			queryString += " AND " + conditions[i]
		}
	}

	return queryString
}
