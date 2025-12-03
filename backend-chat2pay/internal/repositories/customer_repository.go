package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"database/sql"
	"time"

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
	query := `
		INSERT INTO customers (name, email, phone, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.WithContext(ctx).Raw(query,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.PasswordHash,
		now,
		now,
	).Scan(customer).Error

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Customer, error) {
	query := `
		SELECT id, name, email, phone, password_hash, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		var customers []entities.Customer
		err := r.db.WithContext(ctx).Raw(query, limit, offset).Scan(&customers).Error
		return customers, err
	}

	var customers []entities.Customer
	err := r.db.WithContext(ctx).Raw(query).Scan(&customers).Error
	return customers, err
}

func (r *customerRepository) FindOneById(ctx context.Context, id uint64) (*entities.Customer, error) {
	query := `
		SELECT id, name, email, phone, password_hash, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	var customer entities.Customer
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if record was found
	if customer.ID == 0 {
		return nil, nil
	}

	return &customer, nil
}

func (r *customerRepository) FindOneByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	query := `
		SELECT id, name, email, phone, password_hash, created_at, updated_at
		FROM customers
		WHERE email = $1
	`

	var customer entities.Customer
	err := r.db.WithContext(ctx).Raw(query, email).Scan(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if record was found
	if customer.ID == 0 {
		return nil, nil
	}

	return &customer, nil
}

func (r *customerRepository) FindOneByPhone(ctx context.Context, phone string) (*entities.Customer, error) {
	query := `
		SELECT id, name, email, phone, password_hash, created_at, updated_at
		FROM customers
		WHERE phone = $1
	`

	var customer entities.Customer
	err := r.db.WithContext(ctx).Raw(query, phone).Scan(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if record was found
	if customer.ID == 0 {
		return nil, nil
	}

	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *entities.Customer) (*entities.Customer, error) {
	query := `
		UPDATE customers
		SET name = $1, email = $2, phone = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, name, email, phone, password_hash, created_at, updated_at
	`

	now := time.Now()
	err := r.db.WithContext(ctx).Raw(query,
		customer.Name,
		customer.Email,
		customer.Phone,
		now,
		customer.ID,
	).Scan(customer).Error

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) UpdatePassword(ctx context.Context, id uint64, passwordHash string) error {
	query := `UPDATE customers SET password_hash = $1, updated_at = $2 WHERE id = $3`
	return r.db.WithContext(ctx).Exec(query, passwordHash, time.Now(), id).Error
}

func (r *customerRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM customers WHERE id = $1`
	return r.db.WithContext(ctx).Exec(query, id).Error
}

func (r *customerRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM customers`

	var count int64
	err := r.db.WithContext(ctx).Raw(query).Scan(&count).Error
	return count, err
}
