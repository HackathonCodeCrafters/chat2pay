package dto

import (
	"chat2pay/internal/entities"
	"time"
)

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type OrderRequest struct {
	CustomerID     string             `json:"customer_id" validate:"required"`
	MerchantID     string             `json:"merchant_id" validate:"required"`
	OutletID       *string            `json:"outlet_id"`
	Items          []OrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingAmount float64            `json:"shipping_amount" validate:"gte=0"`
	DiscountAmount float64            `json:"discount_amount" validate:"gte=0"`
}

type OrderItemResponse struct {
	ID                  string         `json:"id"`
	ProductID           string         `json:"product_id"`
	Product             *ProductSimple `json:"product,omitempty"`
	ProductNameSnapshot string         `json:"product_name_snapshot"`
	UnitPrice           float64        `json:"unit_price"`
	Quantity            int            `json:"quantity"`
	TotalPrice          float64        `json:"total_price"`
	CreatedAt           time.Time      `json:"created_at"`
}

type ProductSimple struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	SKU  *string `json:"sku,omitempty"`
}

type CustomerSimple struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

type MerchantSimple struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type OrderResponse struct {
	ID             string              `json:"id"`
	OrderNumber    string              `json:"order_number"`
	CustomerID     string              `json:"customer_id"`
	Customer       *CustomerSimple     `json:"customer,omitempty"`
	MerchantID     string              `json:"merchant_id"`
	Merchant       *MerchantSimple     `json:"merchant,omitempty"`
	OutletID       *string             `json:"outlet_id,omitempty"`
	Status         string              `json:"status"`
	SubtotalAmount float64             `json:"subtotal_amount"`
	ShippingAmount float64             `json:"shipping_amount"`
	DiscountAmount float64             `json:"discount_amount"`
	TotalAmount    float64             `json:"total_amount"`
	Items          []OrderItemResponse `json:"items,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type OrderStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=pending paid shipped completed cancelled"`
}

func ToOrderResponse(order *entities.Order) OrderResponse {
	response := OrderResponse{
		ID:             order.ID.String(),
		OrderNumber:    order.OrderNumber,
		CustomerID:     order.CustomerID.String(),
		MerchantID:     order.MerchantID.String(),
		OutletID:       nil,
		Status:         order.Status,
		SubtotalAmount: order.SubtotalAmount,
		ShippingAmount: order.ShippingAmount,
		DiscountAmount: order.DiscountAmount,
		TotalAmount:    order.TotalAmount,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

	if order.Customer != nil {
		response.Customer = &CustomerSimple{
			ID:    order.Customer.ID.String(),
			Name:  order.Customer.Name,
			Email: order.Customer.Email,
			Phone: order.Customer.Phone,
		}
	}

	if order.Merchant != nil {
		response.Merchant = &MerchantSimple{
			ID:    order.Merchant.ID.String(),
			Name:  order.Merchant.Name,
			Email: order.Merchant.Email,
		}
	}

	if len(order.Items) > 0 {
		items := make([]OrderItemResponse, len(order.Items))
		for i, item := range order.Items {
			itemResp := OrderItemResponse{
				ID:                  item.ID.String(),
				ProductID:           item.ProductID.String(),
				ProductNameSnapshot: item.ProductNameSnapshot,
				UnitPrice:           item.UnitPrice,
				Quantity:            item.Quantity,
				TotalPrice:          item.TotalPrice,
				CreatedAt:           item.CreatedAt,
			}

			if item.Product != nil {
				itemResp.Product = &ProductSimple{
					ID:   item.Product.ID.String(),
					Name: item.Product.Name,
					SKU:  item.Product.SKU,
				}
			}

			items[i] = itemResp
		}
		response.Items = items
	}

	return response
}

func ToOrderListResponse(orders []entities.Order, total int64, page, limit int) OrderListResponse {
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = ToOrderResponse(&order)
	}

	return OrderListResponse{
		Orders: orderResponses,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}
}
