-- Откат начальной схемы базы данных

-- Удаляем индексы
DROP INDEX IF EXISTS idx_material_movements_type;
DROP INDEX IF EXISTS idx_material_movements_material;
DROP INDEX IF EXISTS idx_order_items_product;
DROP INDEX IF EXISTS idx_order_items_order;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_partner;
DROP INDEX IF EXISTS idx_sales_history_date;
DROP INDEX IF EXISTS idx_sales_history_product;
DROP INDEX IF EXISTS idx_sales_history_partner;
DROP INDEX IF EXISTS idx_partners_inn;
DROP INDEX IF EXISTS idx_product_materials_material;
DROP INDEX IF EXISTS idx_product_materials_product;
DROP INDEX IF EXISTS idx_materials_material_type;
DROP INDEX IF EXISTS idx_materials_article;
DROP INDEX IF EXISTS idx_products_product_type;
DROP INDEX IF EXISTS idx_products_article;

-- Удаляем таблицы в обратном порядке (с учетом зависимостей)
DROP TABLE IF EXISTS material_movements;
DROP TABLE IF EXISTS material_supplies;
DROP TABLE IF EXISTS suppliers;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS sales_history;
DROP TABLE IF EXISTS partner_sales_points;
DROP TABLE IF EXISTS partners;
DROP TABLE IF EXISTS partner_types;
DROP TABLE IF EXISTS product_materials;
DROP TABLE IF EXISTS materials;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS measurement_units;
DROP TABLE IF EXISTS material_types;
DROP TABLE IF EXISTS product_types; 