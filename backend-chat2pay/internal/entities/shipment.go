package entities

import "time"

type Shipment struct {
	ID                 string     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID            string     `gorm:"not null;uniqueIndex:uq_shipments_order_id" json:"order_id"`
	CourierName        *string    `gorm:"type:varchar(100)" json:"courier_name,omitempty"`
	TrackingNumber     *string    `gorm:"type:varchar(100)" json:"tracking_number,omitempty"`
	Status             string     `gorm:"type:enum('pending','picked_up','in_transit','delivered','returned');not null;default:'pending';index:idx_shipments_status" json:"status"`
	ShippingAddress    *string    `gorm:"type:text" json:"shipping_address,omitempty"`
	ShippingCity       *string    `gorm:"type:varchar(100)" json:"shipping_city,omitempty"`
	ShippingPostalCode *string    `gorm:"type:varchar(20)" json:"shipping_postal_code,omitempty"`
	ShippedAt          *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt        *time.Time `json:"delivered_at,omitempty"`
	CreatedAt          time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Order              Order      `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
