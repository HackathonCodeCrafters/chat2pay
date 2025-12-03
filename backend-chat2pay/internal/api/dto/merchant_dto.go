package dto

import (
	"chat2pay/internal/entities"
	"time"
)

type MerchantRequest struct {
	Name      string `json:"name" validate:"required"`
	LegalName string `json:"legal_name"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}

type MerchantResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	LegalName *string   `json:"legal_name,omitempty"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MerchantListResponse struct {
	Merchants []MerchantResponse `json:"merchants"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

func ToMerchantResponse(merchant *entities.Merchant) MerchantResponse {
	return MerchantResponse{
		ID:        merchant.ID,
		Name:      merchant.Name,
		LegalName: merchant.LegalName,
		Email:     merchant.Email,
		Phone:     merchant.Phone,
		Status:    merchant.Status,
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
	}
}

func ToMerchantListResponse(merchants []entities.Merchant, total int64, page, limit int) MerchantListResponse {
	merchantResponses := make([]MerchantResponse, len(merchants))
	for i, merchant := range merchants {
		merchantResponses[i] = ToMerchantResponse(&merchant)
	}

	return MerchantListResponse{
		Merchants: merchantResponses,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}
}
