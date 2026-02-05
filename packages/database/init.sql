-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Create stores table with geographic data
CREATE TABLE IF NOT EXISTS stores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    barcode VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create prices table
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    store_id INTEGER REFERENCES stores(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'JPY',
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create spatial index on stores location
CREATE INDEX IF NOT EXISTS idx_stores_location ON stores USING GIST(location);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_prices_store_id ON prices(store_id);
CREATE INDEX IF NOT EXISTS idx_prices_product_id ON prices(product_id);
CREATE INDEX IF NOT EXISTS idx_prices_recorded_at ON prices(recorded_at);

-- Insert sample data (Tokyo area stores)
INSERT INTO stores (name, address, location) VALUES
('セブンイレブン 渋谷店', '東京都渋谷区道玄坂1-2-3', ST_GeographyFromText('POINT(139.7007 35.6595)')),
('ファミリーマート 新宿店', '東京都新宿区新宿3-1-1', ST_GeographyFromText('POINT(139.7034 35.6938)')),
('ローソン 池袋店', '東京都豊島区南池袋1-28-1', ST_GeographyFromText('POINT(139.7101 35.7295)')),
('イオン 品川店', '東京都港区高輪4-10-18', ST_GeographyFromText('POINT(139.7394 35.6284)')),
('西友 上野店', '東京都台東区上野4-8-4', ST_GeographyFromText('POINT(139.7744 35.7074)'));

-- Insert sample products
INSERT INTO products (name, category, barcode) VALUES
('コカ・コーラ 500ml', '飲料', '4902102072706'),
('白米 5kg', '食品', '4901010012345'),
('牛乳 1L', '乳製品', '4901050012345'),
('食パン 6枚', 'パン', '4901810012345'),
('卵 10個入り', '生鮮食品', '4901820012345');

-- Insert sample prices
INSERT INTO prices (store_id, product_id, price, recorded_at) VALUES
(1, 1, 120.00, NOW() - INTERVAL '1 day'),
(2, 1, 115.00, NOW() - INTERVAL '1 day'),
(3, 1, 125.00, NOW() - INTERVAL '1 day'),
(1, 2, 1980.00, NOW() - INTERVAL '2 days'),
(2, 2, 1950.00, NOW() - INTERVAL '2 days'),
(4, 2, 1850.00, NOW() - INTERVAL '2 days'),
(1, 3, 198.00, NOW()),
(2, 3, 185.00, NOW()),
(5, 3, 180.00, NOW()),
(1, 4, 128.00, NOW() - INTERVAL '3 hours'),
(3, 4, 135.00, NOW() - INTERVAL '3 hours'),
(5, 4, 118.00, NOW() - INTERVAL '3 hours');
