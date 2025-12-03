package entities

import "time"

type Product struct {
	ID          uint64           `json:"id"`
	MerchantID  uint64           `json:"merchant_id"`
	OutletID    *uint64          `json:"outlet_id,omitempty"`
	CategoryID  *uint64          `json:"category_id,omitempty"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	SKU         *string          `json:"sku,omitempty"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Merchant    *Merchant        `json:"merchant"`
	Outlet      *Outlet          `json:"outlet"`
	Category    *ProductCategory `json:"category"`
	Images      []ProductImage   `json:"images"`
}
