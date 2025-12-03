package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"github.com/jmoiron/sqlx"
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
	DB *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	query := `
		INSERT INTO product (
			merchant_id, outlet_id, category_id, name, description, sku,
			price, stock, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at;
	`

	err := r.DB.QueryRowContext(ctx, query,
		product.MerchantID,
		product.OutletID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Stock,
		product.Status,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return product, err
}

func (r *productRepository) FindAll(ctx context.Context, merchantId uint64, limit, offset int) ([]entities.Product, error) {
	products := []entities.Product{}

	query := `
		SELECT 
			id, merchant_id, outlet_id, category_id, name, description, sku,
			price, stock, status, created_at, updated_at
		FROM product 
		WHERE ($1 = 0 OR merchant_id = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	err := r.DB.SelectContext(ctx, &products, query, merchantId, limit, offset)
	return products, err
}

func (r *productRepository) FindOneById(ctx context.Context, id uint64) (*entities.Product, error) {
	var p entities.Product

	query := `
		SELECT 
			id, merchant_id, outlet_id, category_id, name, description, sku,
			price, stock, status, created_at, updated_at
		FROM product WHERE id = $1 LIMIT 1;
	`

	err := r.DB.GetContext(ctx, &p, query, id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) FindByCategoryId(ctx context.Context, categoryId uint64, limit, offset int) ([]entities.Product, error) {
	products := []entities.Product{}

	query := `
		SELECT 
			id, merchant_id, outlet_id, category_id, name, description, sku,
			price, stock, status, created_at, updated_at
		FROM product 
		WHERE category_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	err := r.DB.SelectContext(ctx, &products, query, categoryId, limit, offset)
	return products, err
}

func (r *productRepository) Update(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	query := `
		UPDATE product
		SET merchant_id=$1, outlet_id=$2, category_id=$3, name=$4,
			description=$5, sku=$6, price=$7, stock=$8, status=$9, 
			updated_at = NOW()
		WHERE id=$10
		RETURNING updated_at;
	`

	err := r.DB.QueryRowContext(ctx, query,
		product.MerchantID,
		product.OutletID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Stock,
		product.Status,
		product.ID,
	).Scan(&product.UpdatedAt)

	return product, err
}

func (r *productRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint64, quantity int) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE product SET stock = stock + $1 WHERE id = $2`,
		quantity, id,
	)
	return err
}

func (r *productRepository) Count(ctx context.Context, merchantId uint64) (int64, error) {
	var count int64

	query := `
		SELECT COUNT(*) FROM product
		WHERE ($1 = 0 OR merchant_id = $1);
	`

	err := r.DB.GetContext(ctx, &count, query, merchantId)
	return count, err
}
