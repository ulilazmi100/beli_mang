-- Creating the order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    item_id UUID REFERENCES items(id),
    quantity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- -- Index on order_items for quick lookups by order_id and merchant_id
-- CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);
-- CREATE INDEX IF NOT EXISTS idx_order_items_item_id ON order_items (item_id);

-- CREATE INDEX IF NOT EXISTS idx_order_items_order_item ON order_items (order_id, item_id);
