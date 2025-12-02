package repositories

import (
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"context"
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
	err := r.db.WithContext(ctx).Create(merchant).Error
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (r *merchantRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Merchant, error) {
	var merchants []entities.Merchant
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&merchants).Error
	if err != nil {
		return nil, err
	}

	return merchants, nil
}

func (r *merchantRepository) FindOneById(ctx context.Context, id uint64) (*entities.Merchant, error) {
	merchant := entities.Merchant{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&merchant).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) FindOneByEmail(ctx context.Context, email string) (*entities.Merchant, error) {
	merchant := entities.Merchant{}
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&merchant).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) Update(ctx context.Context, merchant *entities.Merchant) (*entities.Merchant, error) {
	err := r.db.WithContext(ctx).Save(merchant).Error
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (r *merchantRepository) Delete(ctx context.Context, id uint64) error {
	err := r.db.WithContext(ctx).Delete(&entities.Merchant{}, id).Error
	return err
}

func (r *merchantRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Merchant{}).Count(&count).Error
	return count, err
}
