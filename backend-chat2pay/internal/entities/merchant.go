package entities

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"time"
)

type Status string

const (
	Suspended           Status = "suspended"
	Active              Status = "active"
	PendingVerification Status = "pending_verification"
)

type Merchant struct {
	gorm.Model
	ID        string    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(150);not null" json:"name"`
	LegalName *string   `gorm:"type:varchar(200)" json:"legal_name,omitempty"`
	Email     string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"email"`
	Phone     *string   `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Status    string    `gorm:"type:merchant_status;not null;default:'pending_verification'" json:"status"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Status) Scan(value interface{}) error {
	*p = Status(value.([]byte))
	return nil
}

func (p Status) Value() (driver.Value, error) {
	return string(p), nil
}
