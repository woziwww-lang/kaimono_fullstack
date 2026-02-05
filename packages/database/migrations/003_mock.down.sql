DELETE FROM prices
WHERE store_id IN (SELECT id FROM stores WHERE name LIKE 'Mock Store %')
   OR product_id IN (SELECT id FROM products WHERE name LIKE 'Mock Product %');

DELETE FROM products WHERE name LIKE 'Mock Product %';
DELETE FROM stores WHERE name LIKE 'Mock Store %';
