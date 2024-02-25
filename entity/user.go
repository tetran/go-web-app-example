package entity

import "time"

type UserID int64

type User struct {
	ID        UserID    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
