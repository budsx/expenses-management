package entity

import "time"

type Expense struct {
	ID           int64
	UserID       int64
	AmountIDR    float64
	Description  string
	ReceiptURL   string
	Status       int32
	AutoApproved bool
	SubmittedAt  time.Time
	ProcessedAt  time.Time
}

type ExpenseApproval struct {
	ExpenseID  int64
	ApproverID int64
	Status     int32
	Notes      string
}
