package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error)
	FindAll(ctx context.Context, limit, offset int) ([]entities.Merchant, error)
	FindOneById(ctx context.Context, id uint64) (*entities.Merchant, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.Merchant, error)
	Update(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error)
	Delete(ctx context.Context, id uint64) error
	Count(ctx context.Context) (int64, error)
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepo(db *gorm.DB) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

func (r *merchantRepository) Create(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error) {
	query := `
		INSERT INTO merchants (name, legal_name, email, phone, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.WithContext(ctx).Raw(query,
		merchant.Name,
		merchant.LegalName,
		merchant.Email,
		merchant.Phone,
		merchant.Status,
		now,
		now,
	).Scan(merchant).Error

	if err != nil {
		return nil, err
	}

	return merchant, nil
}

func (r *merchantRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Merchant, error) {
	query := `
		SELECT id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		var merchants []entities.Merchant
		err := r.db.WithContext(ctx).Raw(query, limit, offset).Scan(&merchants).Error
		return merchants, err
	}

	var merchants []entities.Merchant
	err := r.db.WithContext(ctx).Raw(query).Scan(&merchants).Error
	return merchants, err
}

func (r *merchantRepository) FindOneById(ctx context.Context, id uint64) (*entities.Merchant, error) {
	query := `
		SELECT id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`

	var merchant entities.Merchant
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if record was found
	if merchant.ID == 0 {
		return nil, nil
	}

	return &merchant, nil
}

func (r *merchantRepository) FindOneByEmail(ctx context.Context, email string) (*entities.Merchant, error) {
	query := `
		SELECT id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		WHERE email = $1
	`

	var merchant entities.Merchant
	err := r.db.WithContext(ctx).Raw(query, email).Scan(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if record was found
	if merchant.ID == 0 {
		return nil, nil
	}

	return &merchant, nil
}

func (r *merchantRepository) Update(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error) {
	query := `
		UPDATE merchants
		SET name = $1, legal_name = $2, email = $3, phone = $4, status = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, name, legal_name, email, phone, status, created_at, updated_at
	`

	now := time.Now()
	err := r.db.WithContext(ctx).Raw(query,
		merchant.Name,
		merchant.LegalName,
		merchant.Email,
		merchant.Phone,
		merchant.Status,
		now,
		merchant.ID,
	).Scan(merchant).Error

	if err != nil {
		return nil, err
	}

	return merchant, nil
}

func (r *merchantRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM merchants WHERE id = $1`
	return r.db.WithContext(ctx).Exec(query, id).Error
}

func (r *merchantRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM merchants`

	var count int64
	err := r.db.WithContext(ctx).Raw(query).Scan(&count).Error
	return count, err
}
