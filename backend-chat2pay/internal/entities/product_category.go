package entities

import "time"

type ProductCategory struct {
	ID        uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string           `gorm:"type:varchar(150);not null" json:"name"`
	ParentID  *uint64          `gorm:"index:idx_product_categories_parent_id" json:"parent_id,omitempty"`
	CreatedAt time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Parent    *ProductCategory `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
