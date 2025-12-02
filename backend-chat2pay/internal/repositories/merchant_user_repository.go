package repositories

import (
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"context"
	"gorm.io/gorm"
)

type MerchantUserRepository interface {
	Create(ctx context.Context, merchantUser *entities.MerchantUser) (*entities.MerchantUser, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.MerchantUser, error)
	FindOneById(ctx context.Context, id uint64) (*entities.MerchantUser, error)
	FindByMerchantId(ctx context.Context, merchantId uint64) ([]entities.MerchantUser, error)
}

type merchantUserRepository struct {
	db *gorm.DB
}

func NewMerchantUserRepo(db *gorm.DB) MerchantUserRepository {
	return &merchantUserRepository{
		db: db,
	}
}

func (r *merchantUserRepository) Create(ctx context.Context, merchantUser *entities.MerchantUser) (*entities.MerchantUser, error) {
	err := r.db.WithContext(ctx).Create(merchantUser).Error
	if err != nil {
		return nil, err
	}
	return merchantUser, nil
}

func (r *merchantUserRepository) FindOneByEmail(ctx context.Context, email string) (*entities.MerchantUser, error) {
	merchantUser := entities.MerchantUser{}
	err := r.db.WithContext(ctx).
		Preload("Merchant").
		Where("email = ?", email).
		First(&merchantUser).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &merchantUser, nil
}

func (r *merchantUserRepository) FindOneById(ctx context.Context, id uint64) (*entities.MerchantUser, error) {
	merchantUser := entities.MerchantUser{}
	err := r.db.WithContext(ctx).
		Preload("Merchant").
		Where("id = ?", id).
		First(&merchantUser).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &merchantUser, nil
}

func (r *merchantUserRepository) FindByMerchantId(ctx context.Context, merchantId uint64) ([]entities.MerchantUser, error) {
	var merchantUsers []entities.MerchantUser
	err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantId).
		Find(&merchantUsers).Error

	if err != nil {
		return nil, err
	}

	return merchantUsers, nil
}
