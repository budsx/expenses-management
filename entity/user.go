package entity

import "time"

type User struct {
	ID        int64
	Email     string
	Name      string
	Role      string
	CreatedAt time.Time
}
