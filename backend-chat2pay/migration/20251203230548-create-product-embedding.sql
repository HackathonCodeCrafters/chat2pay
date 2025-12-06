
-- +migrate Up
CREATE TABLE product_embedding (
                                   id uuid PRIMARY KEY,
                                   product_id uuid REFERENCES product(id),
                                   content TEXT NOT NULL,
                                   embedding vector(1024)
);

-- +migrate Down
DROP TABLE IF EXISTS product_embedding;
