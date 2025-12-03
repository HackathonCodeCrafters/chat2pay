
-- +migrate Up
CREATE TABLE product_embedding (
                                   id SERIAL PRIMARY KEY,
                                   product_id INT REFERENCES product(id),
                                   content TEXT NOT NULL,
                                   embedding vector(1536)
);

-- +migrate Down
DROP TABLE IF EXISTS product_embedding;
