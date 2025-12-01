package entities

import "time"

type Conversation struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID *uint64   `gorm:"index:idx_conversations_customer_id" json:"customer_id,omitempty"`
	MerchantID uint64    `gorm:"not null;index:idx_conversations_merchant_id" json:"merchant_id"`
	Status     string    `gorm:"type:enum('open','closed');not null;default:'open'" json:"status"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Customer   *Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Merchant   Merchant  `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
