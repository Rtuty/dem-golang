-- Откат тестовых данных

-- Удаляем в обратном порядке (соблюдая зависимости)
DELETE FROM material_movements WHERE id > 0;
DELETE FROM material_supplies WHERE id > 0;
DELETE FROM order_items WHERE id > 0;
DELETE FROM orders WHERE id > 0;
DELETE FROM sales_history WHERE id > 0;
DELETE FROM partner_sales_points WHERE id > 0;
DELETE FROM partners WHERE id > 0;
DELETE FROM employees WHERE id > 0;
DELETE FROM suppliers WHERE id > 0;
DELETE FROM product_materials WHERE id > 0;
DELETE FROM products WHERE id > 0;
DELETE FROM materials WHERE id > 0;

-- Сбрасываем последовательности для автоинкрементных полей
ALTER SEQUENCE material_movements_id_seq RESTART WITH 1;
ALTER SEQUENCE material_supplies_id_seq RESTART WITH 1;
ALTER SEQUENCE order_items_id_seq RESTART WITH 1;
ALTER SEQUENCE orders_id_seq RESTART WITH 1;
ALTER SEQUENCE sales_history_id_seq RESTART WITH 1;
ALTER SEQUENCE partner_sales_points_id_seq RESTART WITH 1;
ALTER SEQUENCE partners_id_seq RESTART WITH 1;
ALTER SEQUENCE employees_id_seq RESTART WITH 1;
ALTER SEQUENCE suppliers_id_seq RESTART WITH 1;
ALTER SEQUENCE product_materials_id_seq RESTART WITH 1;
ALTER SEQUENCE products_id_seq RESTART WITH 1;
ALTER SEQUENCE materials_id_seq RESTART WITH 1; 