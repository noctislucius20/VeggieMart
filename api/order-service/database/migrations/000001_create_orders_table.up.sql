CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_code VARCHAR(64) NOT NULL,
    buyer_id BIGINT NOT NULL,
    order_date DATE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    shipping_type VARCHAR(20) NOT NULL DEFAULT 'PICKUP',
    shipping_fee DECIMAL(10, 2) NOT NULL DEFAULT 0,
    order_time TIME NULL,
    remarks TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE UNIQUE INDEX orders_order_code_deletedat_key
ON orders (order_code)
WHERE deleted_at IS NULL;