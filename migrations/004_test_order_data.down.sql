-- Откат тестовых данных

DELETE FROM sales_history WHERE id = 1;
DELETE FROM order_items WHERE id = 1;
DELETE FROM orders WHERE id = 1;
DELETE FROM partners WHERE id = 1; 