package entities

import "time"

type Merchant struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(150);not null" json:"name"`
	LegalName *string   `gorm:"type:varchar(200)" json:"legal_name,omitempty"`
	Email     string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"email"`
	Phone     *string   `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Status    string    `gorm:"type:enum('active','suspended','pending_verification');not null;default:'pending_verification'" json:"status"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
