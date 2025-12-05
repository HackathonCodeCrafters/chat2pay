package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"github.com/jmoiron/sqlx"
)

type MerchantRepository interface {
	Create(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error)
	FindAll(ctx context.Context, limit, offset int) ([]entities.Merchant, error)
	FindOneById(ctx context.Context, id string) (*entities.Merchant, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.Merchant, error)
	Update(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

type merchantRepository struct {
	DB *sqlx.DB
}

func NewMerchantRepo(db *sqlx.DB) MerchantRepository {
	return &merchantRepository{
		DB: db,
	}
}

func (r *merchantRepository) Create(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error) {
	query := `
		INSERT INTO merchants (name, legal_name, email, phone, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`

	err := r.DB.QueryRowContext(ctx, query,
		merchant.Name,
		merchant.LegalName,
		merchant.Email,
		merchant.Phone,
		merchant.Status,
	).Scan(&merchant.ID, &merchant.CreatedAt, &merchant.UpdatedAt)

	return merchant, err
}

func (r *merchantRepository) FindOneById(ctx context.Context, id string) (*entities.Merchant, error) {
	var merchant entities.Merchant

	query := `
		SELECT 
			id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		WHERE id = $1
		LIMIT 1;
	`

	err := r.DB.GetContext(ctx, &merchant, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Merchant, error) {
	var merchants []entities.Merchant

	query := `
		SELECT 
			id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`

	err := r.DB.SelectContext(ctx, &merchants, query, limit, offset)
	return merchants, err
}

func (r *merchantRepository) FindOneByEmail(ctx context.Context, email string) (*entities.Merchant, error) {
	var merchant entities.Merchant

	query := `
		SELECT 
			id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		WHERE email = $1
		LIMIT 1;
	`

	err := r.DB.GetContext(ctx, &merchant, query, email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) Update(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error) {
	query := `
		UPDATE merchants
		SET name=$1, legal_name=$2, email=$3, phone=$4, status=$5, updated_at=NOW()
		WHERE id = $6
		RETURNING updated_at;
	`

	err := r.DB.QueryRowContext(ctx, query,
		merchant.Name,
		merchant.LegalName,
		merchant.Email,
		merchant.Phone,
		merchant.Status,
		merchant.ID,
	).Scan(&merchant.UpdatedAt)

	return merchant, err
}

func (r *merchantRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM merchants WHERE id = $1`, id)
	return err
}

func (r *merchantRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.DB.GetContext(ctx, &count, `SELECT COUNT(*) FROM merchants`)
	return count, err
}
