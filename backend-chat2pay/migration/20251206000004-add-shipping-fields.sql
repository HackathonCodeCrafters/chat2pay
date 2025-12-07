-- +migrate Up

-- Add shipping fields to product
ALTER TABLE product ADD COLUMN IF NOT EXISTS weight INT DEFAULT 0;
ALTER TABLE product ADD COLUMN IF NOT EXISTS length INT DEFAULT 0;
ALTER TABLE product ADD COLUMN IF NOT EXISTS width INT DEFAULT 0;
ALTER TABLE product ADD COLUMN IF NOT EXISTS height INT DEFAULT 0;

-- Add address fields to merchants
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS city_id VARCHAR(10);
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS city_name VARCHAR(100);
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS province_id VARCHAR(10);
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS province_name VARCHAR(100);
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS postal_code VARCHAR(10);

-- Add address fields to customers
ALTER TABLE customers ADD COLUMN IF NOT EXISTS city_id VARCHAR(10);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS city_name VARCHAR(100);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS province_id VARCHAR(10);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS province_name VARCHAR(100);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS postal_code VARCHAR(10);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,
    shipping_cost DECIMAL(15,2) NOT NULL DEFAULT 0,
    total DECIMAL(15,2) NOT NULL DEFAULT 0,
    courier VARCHAR(50),
    courier_service VARCHAR(100),
    shipping_etd VARCHAR(50),
    tracking_number VARCHAR(100),
    shipping_address TEXT,
    shipping_city VARCHAR(100),
    shipping_province VARCHAR(100),
    shipping_postal_code VARCHAR(10),
    payment_method VARCHAR(50),
    payment_status VARCHAR(50) DEFAULT 'pending',
    payment_token VARCHAR(255),
    payment_url TEXT,
    paid_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES product(id),
    product_name VARCHAR(200) NOT NULL,
    product_price DECIMAL(15,2) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    subtotal DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_merchant_id ON orders(merchant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

-- +migrate Down
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

ALTER TABLE customers DROP COLUMN IF EXISTS city_id;
ALTER TABLE customers DROP COLUMN IF EXISTS city_name;
ALTER TABLE customers DROP COLUMN IF EXISTS province_id;
ALTER TABLE customers DROP COLUMN IF EXISTS province_name;
ALTER TABLE customers DROP COLUMN IF EXISTS address;
ALTER TABLE customers DROP COLUMN IF EXISTS postal_code;

ALTER TABLE merchants DROP COLUMN IF EXISTS city_id;
ALTER TABLE merchants DROP COLUMN IF EXISTS city_name;
ALTER TABLE merchants DROP COLUMN IF EXISTS province_id;
ALTER TABLE merchants DROP COLUMN IF EXISTS province_name;
ALTER TABLE merchants DROP COLUMN IF EXISTS address;
ALTER TABLE merchants DROP COLUMN IF EXISTS postal_code;

ALTER TABLE product DROP COLUMN IF EXISTS weight;
ALTER TABLE product DROP COLUMN IF EXISTS length;
ALTER TABLE product DROP COLUMN IF EXISTS width;
ALTER TABLE product DROP COLUMN IF EXISTS height;
