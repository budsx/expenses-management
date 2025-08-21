package entity

import "time"

type Approval struct {
	ID         int64
	ApproverID int64
	ExpenseID  int64
	Status     string
	Notes      string
	CreatedAt  time.Time
}
