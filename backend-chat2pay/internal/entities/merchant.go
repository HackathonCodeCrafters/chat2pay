package entities

import "time"

type Status string

const (
	Suspended           Status = "suspended"
	Active              Status = "active"
	PendingVerification Status = "pending_verification"
)

type Merchant struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	LegalName *string   `db:"legal_name" json:"legal_name,omitempty"`
	Email     string    `db:"email" json:"email"`
	Phone     *string   `db:"phone" json:"phone,omitempty"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
