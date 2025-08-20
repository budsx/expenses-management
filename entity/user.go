package entity

import "time"

type User struct {
	ID           int64
	Email        string
	Name         string
	Role         int // 1=admin, 2=manager, 3=employee
	PasswordHash string
	CreatedAt    time.Time
}
