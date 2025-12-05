-- +migrate Up
CREATE TYPE product_status_enum AS ENUM ('active', 'inactive', 'archived');
CREATE TABLE IF NOT EXISTS product (
                                        id uuid PRIMARY KEY,
                                        merchant_id uuid NOT NULL,
                                        outlet_id uuid NULL,
                                        category_id uuid NULL,
                                        name VARCHAR(200) NOT NULL,
    description TEXT NULL,
    sku VARCHAR(100) NULL,
    price DECIMAL(15,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    status product_status_enum NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS uq_product_merchant_sku
    ON product (merchant_id, sku);

CREATE INDEX IF NOT EXISTS idx_product_merchant_id ON product (merchant_id);
CREATE INDEX IF NOT EXISTS idx_product_outlet_id ON product (outlet_id);
CREATE INDEX IF NOT EXISTS idx_product_category_id ON product (category_id);
CREATE INDEX IF NOT EXISTS idx_product_status ON product (status);



-- +migrate Down

DROP TABLE IF EXISTS product;

DROP TYPE IF EXISTS product_status_enum;
