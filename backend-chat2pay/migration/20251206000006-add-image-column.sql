-- +migrate Up

ALTER TABLE product ADD COLUMN IF NOT EXISTS image TEXT;

-- +migrate Down
ALTER TABLE product DROP COLUMN IF EXISTS image;
