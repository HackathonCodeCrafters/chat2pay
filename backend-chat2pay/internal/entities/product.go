package entities

import (
	"time"
)

type (
	Product struct {
		ID          string           `json:"id"`
		MerchantID  string           `json:"merchant_id"`
		OutletID    *string          `json:"outlet_id,omitempty"`
		CategoryID  *string          `json:"category_id,omitempty"`
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

	ProductEmbedding struct {
		ID         string    `json:"id"`
		ProductId  string    `json:"product_id"`
		Content    string    `json:"content"`
		Embedding  []float32 `json:"embedding" pg:"type:vector(3)"`
		Similarity float64   `json:"distance"`
	}
)
