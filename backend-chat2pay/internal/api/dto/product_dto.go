package dto

import (
	"chat2pay/internal/entities"
	"time"
)

type AskProduct struct {
	Prompt string `json:"prompt"`
}

type ProductRequest struct {
	MerchantID  string  `json:"merchant_id" validate:"required"`
	OutletID    *string `json:"outlet_id"`
	CategoryID  *string `json:"category_id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	Status      string  `json:"status"`
}

type ProductResponse struct {
	ID          string                 `json:"id"`
	MerchantID  string                 `json:"merchant_id"`
	OutletID    *string                `json:"outlet_id,omitempty"`
	CategoryID  *string                `json:"category_id,omitempty"`
	Category    *ProductCategorySimple `json:"category,omitempty"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	SKU         *string                `json:"sku,omitempty"`
	Price       float64                `json:"price"`
	Stock       int                    `json:"stock"`
	Status      string                 `json:"status"`
	Images      []ProductImageResponse `json:"images,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ProductCategorySimple struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductImageResponse struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

func ToProductResponse(product *entities.Product) ProductResponse {
	response := ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		OutletID:    product.OutletID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		Stock:       product.Stock,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	if product.Category != nil {
		response.Category = &ProductCategorySimple{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	if len(product.Images) > 0 {
		images := make([]ProductImageResponse, len(product.Images))
		for i, img := range product.Images {
			images[i] = ProductImageResponse{
				ID:        img.ID,
				ImageURL:  img.ImageURL,
				IsPrimary: img.IsPrimary,
			}
		}
		response.Images = images
	}

	return response
}

func ToProductListResponse(products []entities.Product, total int64, page, limit int) ProductListResponse {
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ToProductResponse(&product)
	}

	return ProductListResponse{
		Products: productResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}
}
