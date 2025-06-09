-- Тестовые данные для проверки ограничений foreign key при удалении продукции

-- Добавляем партнера для тестирования
INSERT INTO partners (id, partner_type_id, company_name, legal_address, inn, director_name, phone, email)
VALUES (1, 1, 'ТестоваяКомпания ООО', 'г. Москва, ул. Тестовая, д. 1', '1234567890', 'Иванов И.И.', '+7(999)123-45-67', 'test@example.com')
ON CONFLICT (id) DO UPDATE SET 
    company_name = EXCLUDED.company_name,
    legal_address = EXCLUDED.legal_address,
    director_name = EXCLUDED.director_name;

-- Добавляем заказ
INSERT INTO orders (id, partner_id, status, total_amount)
VALUES (1, 1, 'created', 1000.00)
ON CONFLICT (id) DO UPDATE SET 
    status = EXCLUDED.status,
    total_amount = EXCLUDED.total_amount;

-- Добавляем позицию заказа, которая ссылается на продукцию ID=1
INSERT INTO order_items (id, order_id, product_id, quantity, unit_price, total_price)
VALUES (1, 1, 1, 5, 200.00, 1000.00)
ON CONFLICT (id) DO UPDATE SET 
    quantity = EXCLUDED.quantity,
    unit_price = EXCLUDED.unit_price,
    total_price = EXCLUDED.total_price;

-- Добавляем запись в историю продаж для проверки второго ограничения
INSERT INTO sales_history (id, partner_id, product_id, quantity, unit_price, total_amount, sale_date)
VALUES (1, 1, 2, 3, 150.00, 450.00, CURRENT_DATE)
ON CONFLICT (id) DO UPDATE SET 
    quantity = EXCLUDED.quantity,
    unit_price = EXCLUDED.unit_price,
    total_amount = EXCLUDED.total_amount; 