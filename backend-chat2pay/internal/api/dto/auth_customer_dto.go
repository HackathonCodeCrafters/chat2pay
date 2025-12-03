package dto

import "chat2pay/internal/entities"

type CustomerRegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=6"`
}

type CustomerLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CustomerAuthResponse struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Role        string  `json:"role"`
	AccessToken string  `json:"access_token"`
}

func ToCustomerAuthResponse(customer *entities.Customer, token string) CustomerAuthResponse {
	return CustomerAuthResponse{
		ID:          customer.ID,
		Name:        customer.Name,
		Email:       customer.Email,
		Phone:       customer.Phone,
		Role:        "customer",
		AccessToken: token,
	}
}
