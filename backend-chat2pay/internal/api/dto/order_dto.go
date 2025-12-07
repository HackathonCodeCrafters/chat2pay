package dto

import (
	"chat2pay/internal/entities"
	"time"
)

type CreateOrderRequest struct {
	Items            []OrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingAddress  string             `json:"shipping_address" validate:"required"`
	ShippingCity     string             `json:"shipping_city" validate:"required"`
	ShippingCityID   string             `json:"shipping_city_id" validate:"required"`
	ShippingProvince string             `json:"shipping_province" validate:"required"`
	ShippingPostalCode string           `json:"shipping_postal_code"`
	Courier          string             `json:"courier" validate:"required"`
	CourierService   string             `json:"courier_service" validate:"required"`
	ShippingCost     float64            `json:"shipping_cost" validate:"required"`
	ShippingEtd      string             `json:"shipping_etd"`
	Notes            string             `json:"notes"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status         string `json:"status"`
	TrackingNumber string `json:"tracking_number"`
}

type OrderResponse struct {
	ID                 string              `json:"id"`
	CustomerID         string              `json:"customer_id"`
	MerchantID         string              `json:"merchant_id"`
	Status             string              `json:"status"`
	Subtotal           float64             `json:"subtotal"`
	ShippingCost       float64             `json:"shipping_cost"`
	Total              float64             `json:"total"`
	Courier            *string             `json:"courier,omitempty"`
	CourierService     *string             `json:"courier_service,omitempty"`
	ShippingEtd        *string             `json:"shipping_etd,omitempty"`
	TrackingNumber     *string             `json:"tracking_number,omitempty"`
	ShippingAddress    *string             `json:"shipping_address,omitempty"`
	ShippingCity       *string             `json:"shipping_city,omitempty"`
	ShippingProvince   *string             `json:"shipping_province,omitempty"`
	ShippingPostalCode *string             `json:"shipping_postal_code,omitempty"`
	PaymentMethod      *string             `json:"payment_method,omitempty"`
	PaymentStatus      string              `json:"payment_status"`
	PaymentURL         *string             `json:"payment_url,omitempty"`
	PaidAt             *time.Time          `json:"paid_at,omitempty"`
	Notes              *string             `json:"notes,omitempty"`
	Items              []OrderItemResponse `json:"items,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

func ToOrderResponse(order *entities.Order) OrderResponse {
	resp := OrderResponse{
		ID:                 order.ID,
		CustomerID:         order.CustomerID,
		MerchantID:         order.MerchantID,
		Status:             order.Status,
		Subtotal:           order.Subtotal,
		ShippingCost:       order.ShippingCost,
		Total:              order.Total,
		Courier:            order.Courier,
		CourierService:     order.CourierService,
		ShippingEtd:        order.ShippingEtd,
		TrackingNumber:     order.TrackingNumber,
		ShippingAddress:    order.ShippingAddress,
		ShippingCity:       order.ShippingCity,
		ShippingProvince:   order.ShippingProvince,
		ShippingPostalCode: order.ShippingPostalCode,
		PaymentMethod:      order.PaymentMethod,
		PaymentStatus:      order.PaymentStatus,
		PaymentURL:         order.PaymentURL,
		PaidAt:             order.PaidAt,
		Notes:              order.Notes,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}

	if order.Items != nil {
		resp.Items = make([]OrderItemResponse, len(order.Items))
		for i, item := range order.Items {
			resp.Items[i] = OrderItemResponse{
				ID:           item.ID,
				ProductID:    item.ProductID,
				ProductName:  item.ProductName,
				ProductPrice: item.ProductPrice,
				Quantity:     item.Quantity,
				Subtotal:     item.Subtotal,
			}
		}
	}

	return resp
}

func ToOrderListResponse(orders []entities.Order, total int64, page, limit int) OrderListResponse {
	resp := OrderListResponse{
		Orders: make([]OrderResponse, len(orders)),
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	for i, order := range orders {
		resp.Orders[i] = ToOrderResponse(&order)
	}

	return resp
}
