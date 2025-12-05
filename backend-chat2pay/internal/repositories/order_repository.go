package repositories

import (
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"context"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) (*entities.Order, error)
	CreateWithItems(ctx context.Context, order *entities.Order, items []entities.OrderItem) (*entities.Order, error)
	FindAll(ctx context.Context, merchantId, customerId string, limit, offset int) ([]entities.Order, error)
	FindOneById(ctx context.Context, id string) (*entities.Order, error)
	FindOneByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) (*entities.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, merchantId, customerId string) (int64, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) CreateWithItems(ctx context.Context, order *entities.Order, items []entities.OrderItem) (*entities.Order, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) FindAll(ctx context.Context, merchantId, customerId string, limit, offset int) ([]entities.Order, error) {
	var orders []entities.Order
	query := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Merchant").
		Preload("Items").
		Preload("Items.Product")

	if merchantId != "" {
		query = query.Where("merchant_id = ?", merchantId)
	}

	if customerId != "" {
		query = query.Where("customer_id = ?", customerId)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindOneById(ctx context.Context, id string) (*entities.Order, error) {
	order := entities.Order{}
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Merchant").
		Preload("Items").
		Preload("Items.Product").
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) FindOneByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error) {
	order := entities.Order{}
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Merchant").
		Preload("Items").
		Preload("Items.Product").
		Where("order_number = ?", orderNumber).
		First(&order).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	err := r.db.WithContext(ctx).Save(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("id = ?", id).
		Update("status", status).
		Error
	return err
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Delete(&entities.Order{}, id).Error
	return err
}

func (r *orderRepository) Count(ctx context.Context, merchantId, customerId string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Order{})

	if merchantId != "" {
		query = query.Where("merchant_id = ?", merchantId)
	}

	if customerId != "" {
		query = query.Where("customer_id = ?", customerId)
	}

	err := query.Count(&count).Error
	return count, err
}
