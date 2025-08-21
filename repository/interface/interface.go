package _interface

import (
	"context"

	"github.com/budsx/expenses-management/entity"
)

type PaymentProcessor interface {
	ProcessPayment(context.Context, *entity.PaymentProcessorRequest) (*entity.PaymentProcessorResponse, error)
}

type UserRepository interface {
	GetUserWithPassword(context.Context, string) (*entity.User, error)
}

type ExpensesRepository interface {
	WriteExpense(context.Context, *entity.Expense) (int64, error)
	ApprovalExpense(context.Context, *entity.ExpenseApproval) error
	UpdateExpenseStatus(context.Context, int64, int32) error
	GetExpenseByID(context.Context, int64) (*entity.Expense, error)
	GetExpensesWithPagination(context.Context, *entity.ExpenseListQuery) ([]*entity.Expense, int64, error)
	WriteAuditLog(context.Context, *entity.AuditLog) error
	PingContext(context.Context) error
}
