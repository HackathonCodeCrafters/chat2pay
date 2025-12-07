package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) (*entities.Product, error)
	FindAll(ctx context.Context, merchantId string, limit, offset int) ([]entities.Product, error)
	FindByIDs(ctx context.Context, ids []string) ([]entities.Product, error)
	FindOneById(ctx context.Context, id string) (*entities.Product, error)
	FindByCategoryId(ctx context.Context, categoryId string, limit, offset int) ([]entities.Product, error)
	Update(ctx context.Context, product *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, quantity int) error
	Count(ctx context.Context, merchantId string) (int64, error)

	CreateProductEmbedding(ctx context.Context, embedding *entities.ProductEmbedding) error
	GetProductEmbedding(ctx context.Context, vector []float32) (*entities.ProductEmbedding, error)
	GetProductEmbeddingList(ctx context.Context, vector []float32) ([]entities.ProductEmbedding, error)
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
		    id, merchant_id, outlet_id, category_id, name, description, sku, price, stock, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at, updated_at;
	`

	err := r.DB.QueryRowContext(ctx, query,
		uuid.New().String(),
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

func (r *productRepository) FindAll(ctx context.Context, merchantId string, limit, offset int) ([]entities.Product, error) {
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

func (r *productRepository) FindByIDs(ctx context.Context, ids []string) ([]entities.Product, error) {
	products := []entities.Product{}

	query := `
		SELECT 
			id, merchant_id, outlet_id, category_id, name, description, sku,
			price, stock, status, created_at, updated_at
		FROM product 
		WHERE id = ANY($1)
		ORDER BY created_at DESC;
	`

	row, err := r.DB.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	for row.Next() {
		p := entities.Product{}
		if err = row.Scan(&p.ID,
			&p.MerchantID,
			&p.OutletID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.SKU,
			&p.Price,
			&p.Stock,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt); err != nil {
			return nil, err
		}

		products = append(products, p)
	}
	return products, err
}

func (r *productRepository) FindOneById(ctx context.Context, id string) (*entities.Product, error) {
	var p entities.Product

	query := `
		SELECT 
			id, merchant_id, outlet_id, category_id, name, description, sku,
			price, stock, status, created_at, updated_at
		FROM product WHERE id = $1 LIMIT 1;
	`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.MerchantID,
		&p.OutletID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.SKU,
		&p.Price,
		&p.Stock,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) FindByCategoryId(ctx context.Context, categoryId string, limit, offset int) ([]entities.Product, error) {
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

func (r *productRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE product SET stock = stock + $1 WHERE id = $2`,
		quantity, id,
	)
	return err
}

func (r *productRepository) Count(ctx context.Context, merchantId string) (int64, error) {
	var count int64

	query := `
		SELECT COUNT(*) FROM product
		WHERE ($1 = 0 OR merchant_id = $1);
	`

	err := r.DB.GetContext(ctx, &count, query, merchantId)
	return count, err
}

func (r *productRepository) GetProductEmbedding(ctx context.Context, vector []float32) (*entities.ProductEmbedding, error) {
	query := fmt.Sprintf(`SELECT id, product_id, embedding <-> $1 AS distance
		FROM product_embedding
		ORDER BY distance ASC
		LIMIT 5;
	`)

	embed := &entities.ProductEmbedding{}
	if err := r.DB.QueryRowContext(ctx, query, pgvector.NewVector(vector)).Scan(&embed.ID, &embed.ProductId, &embed.Similarity); err != nil {
		return nil, err
	}

	return embed, nil
}

func (r *productRepository) GetProductEmbeddingList(ctx context.Context, vector []float32) ([]entities.ProductEmbedding, error) {
	embeddingQuery := `
        SELECT 
            pe.id,
            pe.product_id,
            1 - (pe.embedding <=> $1) as similarity_score  -- Convert distance to similarity
        FROM product_embedding pe
        JOIN product p ON pe.product_id = p.id
        WHERE p.status = 'active'::public.product_status_enu
        AND p.stock > 0
        AND 1 - (pe.embedding <=> $1) > $2  -- Minimum similarity threshold (0.3 = 70% similarity)
        ORDER BY similarity_score DESC
        LIMIT $3
    `

	// Similarity threshold: 0.3 means we want at least 70% similarity
	// 1 - (cosine distance) = similarity
	// Example: distance 0.3 → similarity 0.7 (70%)
	rows, err := r.DB.QueryContext(ctx, embeddingQuery,
		pgvector.NewVector(vector),
		0.3, // Minimum 70% similarity
		10,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entities.ProductEmbedding
	for rows.Next() {
		var pe entities.ProductEmbedding
		if err := rows.Scan(&pe.ID, &pe.ProductId, &pe.Similarity); err != nil {
			return nil, err
		}
		results = append(results, pe)
	}

	return results, nil
}
func (r *productRepository) CreateProductEmbedding(ctx context.Context, embedding *entities.ProductEmbedding) error {
	query := `
		INSERT INTO product_embedding (
			id, product_id, content, embedding
		) VALUES ($1,$2,$3,$4);
	`

	_, err := r.DB.ExecContext(ctx, query, uuid.New().String(), embedding.ProductId, embedding.Content, pgvector.NewVector(embedding.Embedding))
	if err != nil {
		return err
	}

	return nil
}
