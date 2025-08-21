package _interface

import (
	"context"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/model"
)

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, payment *model.PaymentProcessorModel) error
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type ExpensesRepository interface {
	WriteExpense(ctx context.Context, expense *entity.Expense) (int64, error)
	ApprovalExpense(ctx context.Context, expenseApproval *entity.ExpenseApproval) (int64, error)
	UpdateExpenseStatus(ctx context.Context, expenseID int64, status int32) error
	GetExpenseByID(ctx context.Context, expenseID int64) (*entity.Expense, error)
	WriteAuditLog(ctx context.Context, auditLog *entity.AuditLog) error
}
