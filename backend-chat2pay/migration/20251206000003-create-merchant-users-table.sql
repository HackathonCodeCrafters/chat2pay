-- +migrate Up
CREATE TABLE IF NOT EXISTS merchant_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'staff',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(merchant_id, email)
);

CREATE INDEX IF NOT EXISTS idx_merchant_users_merchant_id ON merchant_users(merchant_id);
CREATE INDEX IF NOT EXISTS idx_merchant_users_email ON merchant_users(email);

-- +migrate Down
DROP TABLE IF EXISTS merchant_users;
