package entities

import "time"

type OrderItem struct {
	ID                  string    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID             string    `gorm:"not null;index:idx_order_items_order_id" json:"order_id"`
	ProductID           string    `gorm:"not null;index:idx_order_items_product_id" json:"product_id"`
	ProductNameSnapshot string    `gorm:"type:varchar(200);not null" json:"product_name_snapshot"`
	UnitPrice           float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	Qty                 int       `gorm:"not null" json:"qty"`
	TotalPrice          float64   `gorm:"type:decimal(15,2);not null" json:"total_price"`
	CreatedAt           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Order               *Order    `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product             *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
