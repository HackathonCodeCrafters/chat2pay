package entities

import "time"

type Customer struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Email        *string   `gorm:"type:varchar(150);uniqueIndex:uq_customers_email" json:"email,omitempty"`
	Phone        *string   `gorm:"type:varchar(50)" json:"phone,omitempty"`
	PasswordHash *string   `gorm:"type:text" json:"password_hash,omitempty"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
