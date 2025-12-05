package entities

import "time"

type Payment struct {
	ID            string     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID       string     `gorm:"not null;uniqueIndex:uq_payments_order_id" json:"order_id"`
	PaymentMethod string     `gorm:"type:varchar(50);not null" json:"payment_method"`
	Provider      *string    `gorm:"type:varchar(50)" json:"provider,omitempty"`
	Amount        float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status        string     `gorm:"type:enum('pending','paid','failed','refunded');not null;default:'pending';index:idx_payments_status" json:"status"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Order         Order      `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
