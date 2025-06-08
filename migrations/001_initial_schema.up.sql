-- Схема базы данных для системы управления производством обоев "Наш декор"
-- Соответствует 3NF с обеспечением ссылочной целостности

-- Типы продукции
CREATE TABLE product_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    coefficient DECIMAL(10,4) NOT NULL DEFAULT 1.0000, -- Коэффициент для расчета материалов
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Типы материалов
CREATE TABLE material_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    defect_rate DECIMAL(5,4) NOT NULL DEFAULT 0.0000, -- Процент брака (0.0500 = 5%)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Единицы измерения
CREATE TABLE measurement_units (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE, -- м, кг, шт, л и т.д.
    symbol VARCHAR(10) NOT NULL UNIQUE, -- м, кг, шт, л
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Продукция
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    article VARCHAR(50) NOT NULL UNIQUE,
    product_type_id INTEGER NOT NULL REFERENCES product_types(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    image_path VARCHAR(500),
    min_partner_price DECIMAL(10,2) NOT NULL CHECK (min_partner_price >= 0),
    
    -- Размеры упаковки (JSON или отдельные поля)
    package_length DECIMAL(10,3), -- длина (м)
    package_width DECIMAL(10,3),  -- ширина (м)  
    package_height DECIMAL(10,3), -- высота (м)
    
    -- Веса
    weight_without_package DECIMAL(10,3), -- вес без упаковки (кг)
    weight_with_package DECIMAL(10,3),    -- вес с упаковкой (кг)
    
    -- Производственные данные
    quality_certificate_path VARCHAR(500), -- путь к скану сертификата
    standard_number VARCHAR(100),          -- номер стандарта
    production_time_hours DECIMAL(8,2),    -- время изготовления (часы)
    cost_price DECIMAL(10,2),             -- себестоимость
    workshop_number VARCHAR(50),           -- номер цеха
    required_workers INTEGER,              -- количество человек на производстве
    roll_width DECIMAL(10,3) CHECK (roll_width >= 0), -- ширина рулона (м)
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Материалы (сырье)
CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    article VARCHAR(50) NOT NULL UNIQUE,
    material_type_id INTEGER NOT NULL REFERENCES material_types(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    measurement_unit_id INTEGER NOT NULL REFERENCES measurement_units(id) ON DELETE RESTRICT,
    package_quantity DECIMAL(10,3) NOT NULL CHECK (package_quantity > 0), -- количество в упаковке
    cost_per_unit DECIMAL(10,2) NOT NULL CHECK (cost_per_unit >= 0),      -- стоимость за единицу
    stock_quantity DECIMAL(10,3) NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0), -- остаток на складе
    min_stock_quantity DECIMAL(10,3) NOT NULL DEFAULT 0 CHECK (min_stock_quantity >= 0), -- минимальный остаток
    image_path VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Связь продукции с материалами (рецептура)
CREATE TABLE product_materials (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    quantity_per_unit DECIMAL(10,6) NOT NULL CHECK (quantity_per_unit > 0), -- количество материала на единицу продукции
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, material_id)
);

-- Партнеры
CREATE TABLE partner_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE partners (
    id SERIAL PRIMARY KEY,
    partner_type_id INTEGER NOT NULL REFERENCES partner_types(id) ON DELETE RESTRICT,
    company_name VARCHAR(200) NOT NULL,
    legal_address TEXT NOT NULL,
    inn VARCHAR(12) NOT NULL UNIQUE,
    director_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    logo_path VARCHAR(500),
    rating INTEGER DEFAULT 0 CHECK (rating >= 0 AND rating <= 10),
    total_sales DECIMAL(15,2) DEFAULT 0 CHECK (total_sales >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Места продаж партнеров
CREATE TABLE partner_sales_points (
    id SERIAL PRIMARY KEY,
    partner_id INTEGER NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    address TEXT NOT NULL,
    sales_type VARCHAR(50) NOT NULL, -- розница, опт, интернет
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- История продаж (для расчета скидок)
CREATE TABLE sales_history (
    id SERIAL PRIMARY KEY,
    partner_id INTEGER NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    total_amount DECIMAL(15,2) NOT NULL CHECK (total_amount >= 0),
    sale_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Сотрудники
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    birth_date DATE NOT NULL,
    passport_series VARCHAR(4),
    passport_number VARCHAR(6),
    bank_details TEXT,
    has_family BOOLEAN DEFAULT FALSE,
    health_status VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Заявки от партнеров
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    partner_id INTEGER NOT NULL REFERENCES partners(id) ON DELETE RESTRICT,
    manager_id INTEGER REFERENCES employees(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'created', -- created, confirmed, prepaid, in_production, ready, completed, cancelled
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    prepayment_amount DECIMAL(15,2) DEFAULT 0,
    delivery_required BOOLEAN DEFAULT FALSE,
    delivery_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Позиции заявок
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    total_price DECIMAL(15,2) NOT NULL CHECK (total_price >= 0),
    production_deadline DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Поставщики материалов
CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    inn VARCHAR(12) NOT NULL UNIQUE,
    contact_info TEXT,
    rating INTEGER DEFAULT 0 CHECK (rating >= 0 AND rating <= 10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- История поставок материалов
CREATE TABLE material_supplies (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE RESTRICT,
    quantity DECIMAL(10,3) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    total_amount DECIMAL(15,2) NOT NULL CHECK (total_amount >= 0),
    supply_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Движение материалов на складе
CREATE TABLE material_movements (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    movement_type VARCHAR(20) NOT NULL, -- income, consumption, write_off, reserve
    quantity DECIMAL(10,3) NOT NULL,
    remaining_quantity DECIMAL(10,3) NOT NULL CHECK (remaining_quantity >= 0),
    reference_id INTEGER, -- ID связанной записи (заявка, поставка и т.д.)
    reference_type VARCHAR(50), -- order, supply, write_off
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_products_article ON products(article);
CREATE INDEX idx_products_product_type ON products(product_type_id);
CREATE INDEX idx_materials_article ON materials(article);
CREATE INDEX idx_materials_material_type ON materials(material_type_id);
CREATE INDEX idx_product_materials_product ON product_materials(product_id);
CREATE INDEX idx_product_materials_material ON product_materials(material_id);
CREATE INDEX idx_partners_inn ON partners(inn);
CREATE INDEX idx_sales_history_partner ON sales_history(partner_id);
CREATE INDEX idx_sales_history_product ON sales_history(product_id);
CREATE INDEX idx_sales_history_date ON sales_history(sale_date);
CREATE INDEX idx_orders_partner ON orders(partner_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_product ON order_items(product_id);
CREATE INDEX idx_material_movements_material ON material_movements(material_id);
CREATE INDEX idx_material_movements_type ON material_movements(movement_type); 