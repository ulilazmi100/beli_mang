-- Creating the merchants table
CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(35) NOT NULL,
    merchant_category VARCHAR(30) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Creating an index on the merchants table for geolocation searches
CREATE INDEX IF NOT EXISTS idx_merchants_location ON merchants USING gist (
    point(latitude, longitude)
);

-- Creating additional indices for optimization
-- Index on merchants for category and name filtering
CREATE INDEX IF NOT EXISTS idx_merchants_category ON merchants (merchant_category);
CREATE INDEX IF NOT EXISTS idx_merchants_name ON merchants (name);

-- CREATE INDEX IF NOT EXISTS idx_merchants_category_name ON merchants (merchant_category, name);