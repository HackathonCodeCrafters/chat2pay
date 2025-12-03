package dto

import "chat2pay/internal/entities"

type MerchantRegisterRequest struct {
	MerchantName string `json:"merchant_name" validate:"required"`
	LegalName    string `json:"legal_name"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone"`
	Name         string `json:"name" validate:"required"`
	Password     string `json:"password" validate:"required,min=6"`
}

type MerchantLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type MerchantAuthResponse struct {
	ID          uint64           `json:"id"`
	MerchantID  uint64           `json:"merchant_id"`
	Merchant    *MerchantSimple2 `json:"merchant,omitempty"`
	Name        string           `json:"name"`
	Email       string           `json:"email"`
	Role        string           `json:"role"`
	Status      string           `json:"status"`
	AccessToken string           `json:"access_token"`
}

type MerchantSimple2 struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func ToMerchantAuthResponse(merchantUser *entities.MerchantUser, token string) MerchantAuthResponse {
	response := MerchantAuthResponse{
		ID:          merchantUser.ID,
		MerchantID:  merchantUser.MerchantID,
		Name:        merchantUser.Name,
		Email:       merchantUser.Email,
		Role:        string(merchantUser.Role),
		Status:      string(merchantUser.Status),
		AccessToken: token,
	}

	if merchantUser.Merchant != nil {
		response.Merchant = &MerchantSimple2{
			ID:     merchantUser.Merchant.ID,
			Name:   merchantUser.Merchant.Name,
			Email:  merchantUser.Merchant.Email,
			Status: merchantUser.Merchant.Status,
		}
	}

	return response
}
