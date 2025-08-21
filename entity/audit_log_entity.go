package entity

import "time"

type AuditLog struct {
	ID           int64
	ExpenseID    int64
	NewStatus    int32
	StatusBefore int32
	Notes        string
	CreatedAt    time.Time
}
