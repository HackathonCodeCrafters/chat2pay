-- 1.1 merchants
CREATE TABLE merchants (
                           id              BIGSERIAL PRIMARY KEY,
                           name            VARCHAR(150) NOT NULL,
                           legal_name      VARCHAR(200),
                           email           VARCHAR(150) NOT NULL,
                           phone           VARCHAR(50),
                           status          VARCHAR(50) NOT NULL DEFAULT 'pending_verification' CHECK (status IN ('active', 'suspended', 'pending_verification')),
                           created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           CONSTRAINT uq_merchants_email UNIQUE (email)
);

-- 1.2 merchant_users
CREATE TABLE merchant_users (
                                id              BIGSERIAL PRIMARY KEY,
                                merchant_id     BIGINT NOT NULL,
                                name            VARCHAR(150) NOT NULL,
                                email           VARCHAR(150) NOT NULL,
                                password_hash   TEXT NOT NULL,
                                role            VARCHAR(50) NOT NULL DEFAULT 'staff' CHECK (role IN ('owner', 'admin', 'staff')),
                                status          VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
                                created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT uq_merchant_users_merchant_email UNIQUE (merchant_id, email),
                                CONSTRAINT fk_merchant_users_merchant
                                    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
                                        ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_merchant_users_merchant_id ON merchant_users(merchant_id);

-- 1.3 outlets
CREATE TABLE outlets (
                         id              BIGSERIAL PRIMARY KEY,
                         merchant_id     BIGINT NOT NULL,
                         name            VARCHAR(150) NOT NULL,
                         address         TEXT,
                         city            VARCHAR(100),
                         latitude        DECIMAL(10,7),
                         longitude       DECIMAL(10,7),
                         phone           VARCHAR(50),
                         status          VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed')),
                         created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         CONSTRAINT fk_outlets_merchant
                             FOREIGN KEY (merchant_id) REFERENCES merchants(id)
                                 ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_outlets_merchant_id ON outlets(merchant_id);

-- 1.4 product_categories
CREATE TABLE product_categories (
                                    id              BIGSERIAL PRIMARY KEY,
                                    name            VARCHAR(150) NOT NULL,
                                    parent_id       BIGINT NULL,
                                    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    CONSTRAINT fk_product_categories_parent
                                        FOREIGN KEY (parent_id) REFERENCES product_categories(id)
                                            ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE INDEX idx_product_categories_parent_id ON product_categories(parent_id);

-- 1.5 customers
CREATE TABLE customers (
                           id              BIGSERIAL PRIMARY KEY,
                           name            VARCHAR(150) NOT NULL,
                           email           VARCHAR(150),
                           phone           VARCHAR(50),
                           created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           CONSTRAINT uq_customers_email UNIQUE (email)
);

-- ============================================
-- 2. PRODUCTS
-- ============================================

-- 2.1 products
CREATE TABLE products (
                          id              BIGSERIAL PRIMARY KEY,
                          merchant_id     BIGINT NOT NULL,
                          outlet_id       BIGINT NULL,
                          category_id     BIGINT NULL,
                          name            VARCHAR(200) NOT NULL,
                          description     TEXT,
                          sku             VARCHAR(100),
                          price           DECIMAL(15,2) NOT NULL,
                          stock           INT NOT NULL DEFAULT 0,
                          status          VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'out_of_stock')),
                          created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT uq_products_merchant_sku UNIQUE (merchant_id, sku),
                          CONSTRAINT fk_products_merchant
                              FOREIGN KEY (merchant_id) REFERENCES merchants(id)
                                  ON UPDATE CASCADE ON DELETE CASCADE,
                          CONSTRAINT fk_products_outlet
                              FOREIGN KEY (outlet_id) REFERENCES outlets(id)
                                  ON UPDATE CASCADE ON DELETE SET NULL,
                          CONSTRAINT fk_products_category
                              FOREIGN KEY (category_id) REFERENCES product_categories(id)
                                  ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE INDEX idx_products_merchant_id ON products(merchant_id);
CREATE INDEX idx_products_outlet_id ON products(outlet_id);
CREATE INDEX idx_products_category_id ON products(category_id);

-- 2.2 product_images
CREATE TABLE product_images (
                                id              BIGSERIAL PRIMARY KEY,
                                product_id      BIGINT NOT NULL,
                                image_url       TEXT NOT NULL,
                                is_primary      BOOLEAN NOT NULL DEFAULT FALSE,
                                created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT fk_product_images_product
                                    FOREIGN KEY (product_id) REFERENCES products(id)
                                        ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_product_images_product_id ON product_images(product_id);

-- ============================================
-- 3. ORDERS & PAYMENTS
-- ============================================

-- 3.1 orders
CREATE TABLE orders (
                        id                  BIGSERIAL PRIMARY KEY,
                        order_number        VARCHAR(50) NOT NULL,
                        customer_id         BIGINT NOT NULL,
                        merchant_id         BIGINT NOT NULL,
                        outlet_id           BIGINT NULL,
                        status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'shipped', 'completed', 'cancelled')),
                        subtotal_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
                        shipping_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
                        discount_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
                        total_amount        DECIMAL(15,2) NOT NULL DEFAULT 0,
                        created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT uq_orders_order_number UNIQUE (order_number),
                        CONSTRAINT fk_orders_customer
                            FOREIGN KEY (customer_id) REFERENCES customers(id)
                                ON UPDATE CASCADE ON DELETE RESTRICT,
                        CONSTRAINT fk_orders_merchant
                            FOREIGN KEY (merchant_id) REFERENCES merchants(id)
                                ON UPDATE CASCADE ON DELETE RESTRICT,
                        CONSTRAINT fk_orders_outlet
                            FOREIGN KEY (outlet_id) REFERENCES outlets(id)
                                ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_merchant_id ON orders(merchant_id);
CREATE INDEX idx_orders_outlet_id ON orders(outlet_id);
CREATE INDEX idx_orders_status ON orders(status);

-- 3.2 order_items
CREATE TABLE order_items (
                             id                      BIGSERIAL PRIMARY KEY,
                             order_id                BIGINT NOT NULL,
                             product_id              BIGINT NOT NULL,
                             product_name_snapshot   VARCHAR(200) NOT NULL,
                             unit_price              DECIMAL(15,2) NOT NULL,
                             qty                     INT NOT NULL,
                             total_price             DECIMAL(15,2) NOT NULL,
                             created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT fk_order_items_order
                                 FOREIGN KEY (order_id) REFERENCES orders(id)
                                     ON UPDATE CASCADE ON DELETE CASCADE,
                             CONSTRAINT fk_order_items_product
                                 FOREIGN KEY (product_id) REFERENCES products(id)
                                     ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- 3.3 payments
CREATE TABLE payments (
                          id              BIGSERIAL PRIMARY KEY,
                          order_id        BIGINT NOT NULL,
                          payment_method  VARCHAR(50) NOT NULL,  -- bank_transfer, ewallet, cod, dll
                          provider        VARCHAR(50),           -- BCA, OVO, DANA, dll
                          amount          DECIMAL(15,2) NOT NULL,
                          status          VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed', 'refunded')),
                          paid_at         TIMESTAMP NULL,
                          created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT uq_payments_order_id UNIQUE (order_id),
                          CONSTRAINT fk_payments_order
                              FOREIGN KEY (order_id) REFERENCES orders(id)
                                  ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_payments_status ON payments(status);

-- 3.4 shipments
CREATE TABLE shipments (
                           id                  BIGSERIAL PRIMARY KEY,
                           order_id            BIGINT NOT NULL,
                           courier_name        VARCHAR(100),
                           tracking_number     VARCHAR(100),
                           status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'picked_up', 'in_transit', 'delivered', 'returned')),
                           shipping_address    TEXT,
                           shipping_city       VARCHAR(100),
                           shipping_postal_code VARCHAR(20),
                           shipped_at          TIMESTAMP NULL,
                           delivered_at        TIMESTAMP NULL,
                           created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           CONSTRAINT uq_shipments_order_id UNIQUE (order_id),
                           CONSTRAINT fk_shipments_order
                               FOREIGN KEY (order_id) REFERENCES orders(id)
                                   ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_shipments_status ON shipments(status);

-- ============================================
-- 4. CHAT / CONVERSATIONS
-- ============================================

-- 4.1 conversations
CREATE TABLE conversations (
                               id              BIGSERIAL PRIMARY KEY,
                               customer_id     BIGINT NULL,
                               merchant_id     BIGINT NOT NULL,
                               status          VARCHAR(50) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed')),
                               created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               CONSTRAINT fk_conversations_customer
                                   FOREIGN KEY (customer_id) REFERENCES customers(id)
                                       ON UPDATE CASCADE ON DELETE SET NULL,
                               CONSTRAINT fk_conversations_merchant
                                   FOREIGN KEY (merchant_id) REFERENCES merchants(id)
                                       ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_conversations_customer_id ON conversations(customer_id);
CREATE INDEX idx_conversations_merchant_id ON conversations(merchant_id);

-- 4.2 messages
CREATE TABLE messages (
                          id              BIGSERIAL PRIMARY KEY,
                          conversation_id BIGINT NOT NULL,
                          sender_type     VARCHAR(50) NOT NULL CHECK (sender_type IN ('customer', 'merchant_user', 'system')),
                          sender_id       BIGINT NULL,
                          message_text    TEXT NOT NULL,
                          created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_messages_conversation
                              FOREIGN KEY (conversation_id) REFERENCES conversations(id)
                                  ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_type_sender_id ON messages(sender_type, sender_id);

-- ============================================
-- TRIGGERS FOR AUTO-UPDATE updated_at
-- ============================================

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to all tables with updated_at
CREATE TRIGGER update_merchants_updated_at BEFORE UPDATE ON merchants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_merchant_users_updated_at BEFORE UPDATE ON merchant_users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_outlets_updated_at BEFORE UPDATE ON outlets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_product_categories_updated_at BEFORE UPDATE ON product_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_shipments_updated_at BEFORE UPDATE ON shipments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_conversations_updated_at BEFORE UPDATE ON conversations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();