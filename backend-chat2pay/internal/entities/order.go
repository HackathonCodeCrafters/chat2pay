package entities

import "time"

type Order struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber    string    `gorm:"type:varchar(50);not null;uniqueIndex:uq_orders_order_number" json:"order_number"`
	CustomerID     uint64    `gorm:"not null;index:idx_orders_customer_id" json:"customer_id"`
	MerchantID     uint64    `gorm:"not null;index:idx_orders_merchant_id" json:"merchant_id"`
	OutletID       *uint64   `gorm:"index:idx_orders_outlet_id" json:"outlet_id,omitempty"`
	Status         string    `gorm:"type:enum('pending','paid','shipped','completed','cancelled');not null;default:'pending';index:idx_orders_status" json:"status"`
	SubtotalAmount float64   `gorm:"type:decimal(15,2);not null;default:0" json:"subtotal_amount"`
	ShippingAmount float64   `gorm:"type:decimal(15,2);not null;default:0" json:"shipping_amount"`
	DiscountAmount float64   `gorm:"type:decimal(15,2);not null;default:0" json:"discount_amount"`
	TotalAmount    float64   `gorm:"type:decimal(15,2);not null;default:0" json:"total_amount"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Customer       Customer  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Merchant       Merchant  `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Outlet         *Outlet   `gorm:"foreignKey:OutletID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
