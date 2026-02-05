INSERT INTO stores (name, address, location) VALUES
('セブンイレブン 渋谷店', '東京都渋谷区道玄坂1-2-3', ST_GeographyFromText('POINT(139.7007 35.6595)')),
('ファミリーマート 新宿店', '東京都新宿区新宿3-1-1', ST_GeographyFromText('POINT(139.7034 35.6938)')),
('ローソン 池袋店', '東京都豊島区南池袋1-28-1', ST_GeographyFromText('POINT(139.7101 35.7295)')),
('イオン 品川店', '東京都港区高輪4-10-18', ST_GeographyFromText('POINT(139.7394 35.6284)')),
('西友 上野店', '東京都台東区上野4-8-4', ST_GeographyFromText('POINT(139.7744 35.7074)'));

INSERT INTO products (name, category, barcode) VALUES
('コカ・コーラ 500ml', '飲料', '4902102072706'),
('白米 5kg', '食品', '4901010012345'),
('牛乳 1L', '乳製品', '4901050012345'),
('食パン 6枚', 'パン', '4901810012345'),
('卵 10個入り', '生鮮食品', '4901820012345');

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
