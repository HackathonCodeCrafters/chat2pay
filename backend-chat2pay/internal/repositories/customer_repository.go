package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entities.Customer) (*entities.Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]entities.Customer, error)
	FindOneById(ctx context.Context, id string) (*entities.Customer, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.Customer, error)
	FindOneByPhone(ctx context.Context, phone string) (*entities.Customer, error)
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
	Update(ctx context.Context, customer *entities.Customer) (*entities.Customer, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) CustomerRepository {
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
	row := r.db.QueryRowxContext(ctx, query,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.PasswordHash,
		now,
		now,
	)

	err := row.Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Customer, error) {
	var customers []entities.Customer

	if limit > 0 {
		query := `
			SELECT id, name, email, phone, password_hash, created_at, updated_at
			FROM customers
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		err := r.db.SelectContext(ctx, &customers, query, limit, offset)
		return customers, err
	}

	query := `
		SELECT id, name, email, phone, password_hash, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &customers, query)
	return customers, err
}

func (r *customerRepository) FindOneById(ctx context.Context, id string) (*entities.Customer, error) {
	query := `
		SELECT id, name, email, phone, password_hash, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	var customer entities.Customer
	err := r.db.GetContext(ctx, &customer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
	err := r.db.GetContext(ctx, &customer, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
	err := r.db.GetContext(ctx, &customer, query, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
	row := r.db.QueryRowxContext(ctx, query,
		customer.Name,
		customer.Email,
		customer.Phone,
		now,
		customer.ID,
	)

	err := row.StructScan(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	query := `UPDATE customers SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), id)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *customerRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM customers`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}
