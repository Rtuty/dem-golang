-- Откат начальных данных

DELETE FROM partner_types WHERE name IN ('Розничная сеть', 'Оптовая компания', 'Интернет-магазин', 'Строительная компания');

DELETE FROM material_types WHERE name IN ('Основа', 'Покрытие', 'Клеящие материалы', 'Красители', 'Упаковочные материалы');

DELETE FROM product_types WHERE name IN ('Виниловые обои', 'Флизелиновые обои', 'Бумажные обои', 'Текстильные обои');

DELETE FROM measurement_units WHERE symbol IN ('м', 'кг', 'шт', 'л', 'м²', 'рул'); 