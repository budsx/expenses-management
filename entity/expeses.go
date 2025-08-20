package entity

import "time"

type Expense struct {
	ID          int64
	UserID      int64
	AmountIDR   float64
	Description string
	ReceiptURL  string
	Status      string
	SubmittedAt time.Time
	ProcessedAt time.Time
}
