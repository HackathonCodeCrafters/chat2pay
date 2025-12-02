package entities

import "time"

type MerchantUser struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID   uint64    `gorm:"not null;index:idx_merchant_users_merchant_id;uniqueIndex:uq_merchant_users_merchant_email,priority:1" json:"merchant_id"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Email        string    `gorm:"type:varchar(150);not null;uniqueIndex:uq_merchant_users_merchant_email,priority:2" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"password_hash"`
	Role         string    `gorm:"type:enum('owner','admin','staff');not null;default:'staff'" json:"role"`
	Status       string    `gorm:"type:enum('active','inactive');not null;default:'active'" json:"status"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Merchant     *Merchant `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
