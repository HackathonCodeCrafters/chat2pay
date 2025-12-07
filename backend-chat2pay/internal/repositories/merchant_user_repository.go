package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MerchantUserRepository interface {
	Create(ctx context.Context, merchantUser *entities.MerchantUser) (*entities.MerchantUser, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.MerchantUser, error)
	FindOneById(ctx context.Context, id string) (*entities.MerchantUser, error)
	FindByMerchantId(ctx context.Context, merchantId string) ([]entities.MerchantUser, error)
}

type merchantUserRepository struct {
	db *sqlx.DB
}

func NewMerchantUserRepo(db *sqlx.DB) MerchantUserRepository {
	return &merchantUserRepository{
		db: db,
	}
}

func (r *merchantUserRepository) Create(ctx context.Context, merchantUser *entities.MerchantUser) (*entities.MerchantUser, error) {
	query := `
		INSERT INTO merchant_users (id, merchant_id, name, email, password_hash, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	merchantUser.ID = uuid.New().String()
	row := r.db.QueryRowxContext(ctx, query,
		merchantUser.ID,
		merchantUser.MerchantID,
		merchantUser.Name,
		merchantUser.Email,
		merchantUser.PasswordHash,
		merchantUser.Role,
		merchantUser.Status,
		now,
		now,
	)

	err := row.Scan(&merchantUser.ID, &merchantUser.CreatedAt, &merchantUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return merchantUser, nil
}

func (r *merchantUserRepository) FindOneByEmail(ctx context.Context, email string) (*entities.MerchantUser, error) {
	query := `
		SELECT id, merchant_id, name, email, password_hash, role, status, created_at, updated_at
		FROM merchant_users
		WHERE email = $1
	`

	var merchantUser entities.MerchantUser
	err := r.db.GetContext(ctx, &merchantUser, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Load merchant
	merchant, err := r.findMerchantById(ctx, merchantUser.MerchantID)
	if err == nil && merchant != nil {
		merchantUser.Merchant = merchant
	}

	return &merchantUser, nil
}

func (r *merchantUserRepository) FindOneById(ctx context.Context, id string) (*entities.MerchantUser, error) {
	query := `
		SELECT id, merchant_id, name, email, password_hash, role, status, created_at, updated_at
		FROM merchant_users
		WHERE id = $1
	`

	var merchantUser entities.MerchantUser
	err := r.db.GetContext(ctx, &merchantUser, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Load merchant
	merchant, err := r.findMerchantById(ctx, merchantUser.MerchantID)
	if err == nil && merchant != nil {
		merchantUser.Merchant = merchant
	}

	return &merchantUser, nil
}

func (r *merchantUserRepository) FindByMerchantId(ctx context.Context, merchantId string) ([]entities.MerchantUser, error) {
	query := `
		SELECT id, merchant_id, name, email, password_hash, role, status, created_at, updated_at
		FROM merchant_users
		WHERE merchant_id = $1
	`

	var merchantUsers []entities.MerchantUser
	err := r.db.SelectContext(ctx, &merchantUsers, query, merchantId)
	if err != nil {
		return nil, err
	}

	return merchantUsers, nil
}

func (r *merchantUserRepository) findMerchantById(ctx context.Context, id string) (*entities.Merchant, error) {
	query := `
		SELECT id, name, legal_name, email, phone, status, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`

	var merchant entities.Merchant
	err := r.db.GetContext(ctx, &merchant, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &merchant, nil
}
