-- Creating the orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    total_price NUMERIC NOT NULL,
    estimated_delivery_time_in_minutes NUMERIC NOT NULL,
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index on orders for user_id
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);

