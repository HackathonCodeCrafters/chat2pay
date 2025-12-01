package entities

import "time"

type Outlet struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID uint64    `gorm:"not null;index:idx_outlets_merchant_id" json:"merchant_id"`
	Name       string    `gorm:"type:varchar(150);not null" json:"name"`
	Address    *string   `gorm:"type:text" json:"address,omitempty"`
	City       *string   `gorm:"type:varchar(100)" json:"city,omitempty"`
	Latitude   *float64  `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude  *float64  `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	Phone      *string   `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Status     string    `gorm:"type:enum('active','closed');not null;default:'active'" json:"status"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Merchant   Merchant  `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
