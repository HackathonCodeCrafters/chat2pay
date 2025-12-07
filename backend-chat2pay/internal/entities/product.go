package entities

import (
	"time"
)

type (
	Product struct {
		ID          string           `json:"id" db:"id"`
		MerchantID  string           `json:"merchant_id" db:"merchant_id"`
		OutletID    *string          `json:"outlet_id,omitempty" db:"outlet_id"`
		CategoryID  *string          `json:"category_id,omitempty" db:"category_id"`
		Name        string           `json:"name" db:"name"`
		Description *string          `json:"description,omitempty" db:"description"`
		SKU         *string          `json:"sku,omitempty" db:"sku"`
		Price       float64          `json:"price" db:"price"`
		Stock       int              `json:"stock" db:"stock"`
		Status      string           `json:"status" db:"status"`
		Image       *string          `json:"image,omitempty" db:"image"`
		Weight      int              `json:"weight" db:"weight"`
		Length      int              `json:"length" db:"length"`
		Width       int              `json:"width" db:"width"`
		Height      int              `json:"height" db:"height"`
		CreatedAt   time.Time        `json:"created_at" db:"created_at"`
		UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
		Merchant    *Merchant        `json:"merchant" db:"-"`
		Outlet      *Outlet          `json:"outlet" db:"-"`
		Category    *ProductCategory `json:"category" db:"-"`
		Images      []ProductImage   `json:"images" db:"-"`
	}

	ProductEmbedding struct {
		ID         string    `json:"id"`
		ProductId  string    `json:"product_id"`
		Content    string    `json:"content"`
		Embedding  []float32 `json:"embedding" pg:"type:vector(3)"`
		Similarity float64   `json:"distance"`
	}
)
