package repositories

import (
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"context"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entities.Customer) (*entities.Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]entities.Customer, error)
	FindOneById(ctx context.Context, id uint64) (*entities.Customer, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.Customer, error)
	FindOneByPhone(ctx context.Context, phone string) (*entities.Customer, error)
	UpdatePassword(ctx context.Context, id uint64, passwordHash string) error
	Update(ctx context.Context, customer *entities.Customer) (*entities.Customer, error)
	Delete(ctx context.Context, id uint64) error
	Count(ctx context.Context) (int64, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) Create(ctx context.Context, customer *entities.Customer) (*entities.Customer, error) {
	err := r.db.WithContext(ctx).Create(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Customer, error) {
	var customers []entities.Customer
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&customers).Error
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *customerRepository) FindOneById(ctx context.Context, id uint64) (*entities.Customer, error) {
	customer := entities.Customer{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&customer).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindOneByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	customer := entities.Customer{}
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindOneByPhone(ctx context.Context, phone string) (*entities.Customer, error) {
	customer := entities.Customer{}
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&customer).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *entities.Customer) (*entities.Customer, error) {
	err := r.db.WithContext(ctx).Save(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) Delete(ctx context.Context, id uint64) error {
	err := r.db.WithContext(ctx).Delete(&entities.Customer{}, id).Error
	return err
}

func (r *customerRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Customer{}).Count(&count).Error
	return count, err
}

func (r *customerRepository) UpdatePassword(ctx context.Context, id uint64, passwordHash string) error {
	err := r.db.WithContext(ctx).
		Model(&entities.Customer{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).
		Error
	return err
}
