package _interface

import (
	"context"

	"github.com/budsx/expenses-management/entity"
)

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, payment *entity.PaymentProcessorRequest) (*entity.PaymentProcessorResponse, error)
}

type UserRepository interface {
	GetUserWithPassword(ctx context.Context, email string) (*entity.User, error)
}

type ExpensesRepository interface {
	WriteExpense(ctx context.Context, expense *entity.Expense) (int64, error)
	ApprovalExpense(ctx context.Context, expenseApproval *entity.ExpenseApproval) error
	UpdateExpenseStatus(ctx context.Context, expenseID int64, status int32) error
	GetExpenseByID(ctx context.Context, expenseID int64) (*entity.Expense, error)
	WriteAuditLog(ctx context.Context, auditLog *entity.AuditLog) error
	PingContext(ctx context.Context) error
}
