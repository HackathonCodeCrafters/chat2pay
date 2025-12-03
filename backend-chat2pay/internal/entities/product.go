package entities

import "time"

type Product struct {
	ID          uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID  uint64           `gorm:"not null;index:idx_products_merchant_id;uniqueIndex:uq_products_merchant_sku,priority:1" json:"merchant_id"`
	OutletID    *uint64          `gorm:"index:idx_products_outlet_id" json:"outlet_id,omitempty"`
	CategoryID  *uint64          `gorm:"index:idx_products_category_id" json:"category_id,omitempty"`
	Name        string           `gorm:"type:varchar(200);not null" json:"name"`
	Description *string          `gorm:"type:text" json:"description,omitempty"`
	SKU         *string          `gorm:"type:varchar(100);uniqueIndex:uq_products_merchant_sku,priority:2" json:"sku,omitempty"`
	Price       float64          `gorm:"type:decimal(15,2);not null" json:"price"`
	Stock       int              `gorm:"not null;default:0" json:"stock"`
	Status      string           `gorm:"type:enum('active','inactive','out_of_stock');not null;default:'active'" json:"status"`
	CreatedAt   time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Merchant    *Merchant        `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Outlet      *Outlet          `gorm:"foreignKey:OutletID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Category    *ProductCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Images      []ProductImage   `gorm:"foreignKey:ProductID"`
}
