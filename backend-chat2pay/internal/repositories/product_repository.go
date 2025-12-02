package repositories

import (
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"context"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) (*entities.Product, error)
	FindAll(ctx context.Context, merchantId uint64, limit, offset int) ([]entities.Product, error)
	FindOneById(ctx context.Context, id uint64) (*entities.Product, error)
	FindByCategoryId(ctx context.Context, categoryId uint64, limit, offset int) ([]entities.Product, error)
	Update(ctx context.Context, product *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id uint64) error
	UpdateStock(ctx context.Context, id uint64, quantity int) error
	Count(ctx context.Context, merchantId uint64) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	err := r.db.WithContext(ctx).Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) FindAll(ctx context.Context, merchantId uint64, limit, offset int) ([]entities.Product, error) {
	var products []entities.Product
	query := r.db.WithContext(ctx).Preload("Category")

	if merchantId > 0 {
		query = query.Where("merchant_id = ?", merchantId)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) FindOneById(ctx context.Context, id uint64) (*entities.Product, error) {
	product := entities.Product{}
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Where("id = ?", id).
		First(&product).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindByCategoryId(ctx context.Context, categoryId uint64, limit, offset int) ([]entities.Product, error) {
	var products []entities.Product
	query := r.db.WithContext(ctx).
		Preload("Category").
		Where("category_id = ?", categoryId)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	err := r.db.WithContext(ctx).Save(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id uint64) error {
	err := r.db.WithContext(ctx).Delete(&entities.Product{}, id).Error
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint64, quantity int) error {
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).
		Error
	return err
}

func (r *productRepository) Count(ctx context.Context, merchantId uint64) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Product{})

	if merchantId > 0 {
		query = query.Where("merchant_id = ?", merchantId)
	}

	err := query.Count(&count).Error
	return count, err
}
