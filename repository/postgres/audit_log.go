package postgres

import (
	"context"

	"github.com/budsx/expenses-management/entity"
)

func (r *ExpensesRepository) WriteAuditLog(ctx context.Context, auditLog *entity.AuditLog) error {
	query := `
		INSERT INTO audit_logs (expense_id, new_status, status_before, notes, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, auditLog.ExpenseID, auditLog.NewStatus, auditLog.StatusBefore, auditLog.Notes, auditLog.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
