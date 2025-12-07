package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	CreateItem(ctx context.Context, item *entities.OrderItem) error
	FindByID(ctx context.Context, id string) (*entities.Order, error)
	FindByCustomerID(ctx context.Context, customerID string, page, limit int) ([]entities.Order, error)
	FindByMerchantID(ctx context.Context, merchantID string, page, limit int) ([]entities.Order, error)
	CountByCustomerID(ctx context.Context, customerID string) (int64, error)
	CountByMerchantID(ctx context.Context, merchantID string) (int64, error)
	GetOrderItems(ctx context.Context, orderID string) ([]entities.OrderItem, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error
	UpdateTrackingNumber(ctx context.Context, id, trackingNumber string) error
}

type orderRepository struct {
	DB *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{DB: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order) error {
	query := `
		INSERT INTO orders (
			id, customer_id, merchant_id, status, subtotal, shipping_cost, total,
			courier, courier_service, shipping_etd, shipping_address, shipping_city,
			shipping_province, shipping_postal_code, payment_status, notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`
	_, err := r.DB.ExecContext(ctx, query,
		order.ID, order.CustomerID, order.MerchantID, order.Status,
		order.Subtotal, order.ShippingCost, order.Total,
		order.Courier, order.CourierService, order.ShippingEtd,
		order.ShippingAddress, order.ShippingCity, order.ShippingProvince,
		order.ShippingPostalCode, order.PaymentStatus, order.Notes,
	)
	return err
}

func (r *orderRepository) CreateItem(ctx context.Context, item *entities.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, product_name, product_price, quantity, subtotal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.DB.ExecContext(ctx, query,
		item.ID, item.OrderID, item.ProductID, item.ProductName,
		item.ProductPrice, item.Quantity, item.Subtotal,
	)
	return err
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*entities.Order, error) {
	var order entities.Order
	query := `SELECT * FROM orders WHERE id = $1`
	err := r.DB.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByCustomerID(ctx context.Context, customerID string, page, limit int) ([]entities.Order, error) {
	var orders []entities.Order
	offset := (page - 1) * limit
	query := `SELECT * FROM orders WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.DB.SelectContext(ctx, &orders, query, customerID, limit, offset)
	return orders, err
}

func (r *orderRepository) FindByMerchantID(ctx context.Context, merchantID string, page, limit int) ([]entities.Order, error) {
	var orders []entities.Order
	offset := (page - 1) * limit
	query := `SELECT * FROM orders WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.DB.SelectContext(ctx, &orders, query, merchantID, limit, offset)
	return orders, err
}

func (r *orderRepository) CountByCustomerID(ctx context.Context, customerID string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM orders WHERE customer_id = $1`
	err := r.DB.GetContext(ctx, &count, query, customerID)
	return count, err
}

func (r *orderRepository) CountByMerchantID(ctx context.Context, merchantID string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM orders WHERE merchant_id = $1`
	err := r.DB.GetContext(ctx, &count, query, merchantID)
	return count, err
}

func (r *orderRepository) GetOrderItems(ctx context.Context, orderID string) ([]entities.OrderItem, error) {
	var items []entities.OrderItem
	query := `SELECT * FROM order_items WHERE order_id = $1`
	err := r.DB.SelectContext(ctx, &items, query, orderID)
	return items, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, id)
	return err
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error {
	query := `UPDATE orders SET payment_status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, paymentStatus, id)
	return err
}

func (r *orderRepository) UpdateTrackingNumber(ctx context.Context, id, trackingNumber string) error {
	query := `UPDATE orders SET tracking_number = $1, status = 'shipped', updated_at = NOW() WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, trackingNumber, id)
	return err
}
