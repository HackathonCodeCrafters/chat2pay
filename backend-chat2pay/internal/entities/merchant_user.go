package entities

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"time"
)

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
	gorm.Model
	ID           string             `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID   string             `gorm:"not null;index:idx_merchant_users_merchant_id;uniqueIndex:uq_merchant_users_merchant_email,priority:1"`
	Name         string             `gorm:"type:varchar(150);not null"`
	Email        string             `gorm:"type:varchar(150);not null;uniqueIndex:uq_merchant_users_merchant_email,priority:2"`
	PasswordHash string             `gorm:"type:text;not null"`
	Role         MerchantUserRole   `gorm:"type:merchant_user_role;not null;default:'staff'"`
	Status       MerchantUserStatus `gorm:"type:merchant_user_status;not null;default:'active'"`
	CreatedAt    time.Time          `gorm:"autoCreateTime"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime"`
	Merchant     *Merchant          `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (r *MerchantUserRole) Scan(value interface{}) error {
	*r = MerchantUserRole(value.([]byte))
	return nil
}

func (r MerchantUserRole) Value() (driver.Value, error) {
	return string(r), nil
}

func (s *MerchantUserStatus) Scan(value interface{}) error {
	*s = MerchantUserStatus(value.([]byte))
	return nil
}

func (s MerchantUserStatus) Value() (driver.Value, error) {
	return string(s), nil
}
