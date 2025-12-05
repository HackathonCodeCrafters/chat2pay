package dto

import (
	"chat2pay/internal/entities"
	"time"
)

type CustomerRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone"`
}

type CustomerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerListResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

func ToCustomerResponse(customer *entities.Customer) CustomerResponse {
	return CustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

func ToCustomerListResponse(customers []entities.Customer, total int64, page, limit int) CustomerListResponse {
	customerResponses := make([]CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = ToCustomerResponse(&customer)
	}

	return CustomerListResponse{
		Customers: customerResponses,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}
}
