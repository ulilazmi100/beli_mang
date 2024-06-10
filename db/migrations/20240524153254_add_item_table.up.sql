-- Creating the items table
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID,
    name VARCHAR(35) NOT NULL,
    product_category VARCHAR(30) NOT NULL,
    price NUMERIC NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Index on items for merchant_id, category, and name filtering
CREATE INDEX IF NOT EXISTS idx_items_merchant_id ON items (merchant_id);
CREATE INDEX IF NOT EXISTS idx_items_category ON items (product_category);
CREATE INDEX IF NOT EXISTS idx_items_name ON items (name);

CREATE INDEX IF NOT EXISTS idx_items_merchant_id_name ON items (merchant_id, name);
