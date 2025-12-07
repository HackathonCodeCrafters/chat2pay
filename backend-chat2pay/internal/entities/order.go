package entities

import "time"

type Order struct {
	ID                  string       `json:"id" db:"id"`
	CustomerID          string       `json:"customer_id" db:"customer_id"`
	MerchantID          string       `json:"merchant_id" db:"merchant_id"`
	Status              string       `json:"status" db:"status"`
	Subtotal            float64      `json:"subtotal" db:"subtotal"`
	ShippingCost        float64      `json:"shipping_cost" db:"shipping_cost"`
	Total               float64      `json:"total" db:"total"`
	Courier             *string      `json:"courier,omitempty" db:"courier"`
	CourierService      *string      `json:"courier_service,omitempty" db:"courier_service"`
	ShippingEtd         *string      `json:"shipping_etd,omitempty" db:"shipping_etd"`
	TrackingNumber      *string      `json:"tracking_number,omitempty" db:"tracking_number"`
	ShippingAddress     *string      `json:"shipping_address,omitempty" db:"shipping_address"`
	ShippingCity        *string      `json:"shipping_city,omitempty" db:"shipping_city"`
	ShippingProvince    *string      `json:"shipping_province,omitempty" db:"shipping_province"`
	ShippingPostalCode  *string      `json:"shipping_postal_code,omitempty" db:"shipping_postal_code"`
	PaymentMethod       *string      `json:"payment_method,omitempty" db:"payment_method"`
	PaymentStatus       string       `json:"payment_status" db:"payment_status"`
	PaymentToken        *string      `json:"payment_token,omitempty" db:"payment_token"`
	PaymentURL          *string      `json:"payment_url,omitempty" db:"payment_url"`
	PaidAt              *time.Time   `json:"paid_at,omitempty" db:"paid_at"`
	Notes               *string      `json:"notes,omitempty" db:"notes"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" db:"updated_at"`
	Items               []OrderItem  `json:"items,omitempty" db:"-"`
	Customer            *Customer    `json:"customer,omitempty" db:"-"`
	Merchant            *Merchant    `json:"merchant,omitempty" db:"-"`
}

type OrderItem struct {
	ID           string    `json:"id" db:"id"`
	OrderID      string    `json:"order_id" db:"order_id"`
	ProductID    string    `json:"product_id" db:"product_id"`
	ProductName  string    `json:"product_name" db:"product_name"`
	ProductPrice float64   `json:"product_price" db:"product_price"`
	Quantity     int       `json:"quantity" db:"quantity"`
	Subtotal     float64   `json:"subtotal" db:"subtotal"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Product      *Product  `json:"product,omitempty" db:"-"`
}

const (
	OrderStatusPending    = "pending"
	OrderStatusPaid       = "paid"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"

	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"
	PaymentStatusExpired = "expired"
)
