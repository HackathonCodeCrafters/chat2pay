package entities

import "time"

type Customer struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        *string   `db:"email" json:"email,omitempty"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	PasswordHash *string   `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
