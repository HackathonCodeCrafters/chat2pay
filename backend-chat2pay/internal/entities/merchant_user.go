package entities

import "time"

type MerchantUserRole string
type MerchantUserStatus string

const (
	RoleOwner MerchantUserRole = "owner"
	RoleAdmin MerchantUserRole = "admin"
	RoleStaff MerchantUserRole = "staff"

	StatusActive   MerchantUserStatus = "active"
	StatusInactive MerchantUserStatus = "inactive"
)

type MerchantUser struct {
	ID           string    `db:"id" json:"id"`
	MerchantID   string    `db:"merchant_id" json:"merchant_id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Merchant     *Merchant `db:"-" json:"merchant,omitempty"`
}
