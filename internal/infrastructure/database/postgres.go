package database

import (
	"database/sql"
	"fmt"
	"time"

	"wallpaper-system/internal/infrastructure/config"

	_ "github.com/lib/pq"
)

// NewPostgresConnection создает новое подключение к PostgreSQL
func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	// Формируем строку подключения
	dsn := cfg.GetDSN()

	// Открываем соединение
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с БД: %w", err)
	}

	// Настраиваем пул соединений
	setupConnectionPool(db, cfg)

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	return db, nil
}

// setupConnectionPool настраивает пул соединений
func setupConnectionPool(db *sql.DB, cfg config.DatabaseConfig) {
	// Максимальное количество открытых соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	// Максимальное количество неактивных соединений
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// Максимальное время жизни соединения
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Максимальное время простоя соединения
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}

// CreateTables создает таблицы базы данных
func CreateTables(db *sql.DB) error {
	queries := []string{
		createProductTypesTable,
		createMaterialTypesTable,
		createMeasurementUnitsTable,
		createMaterialsTable,
		createProductsTable,
		createProductMaterialsTable,
		createIndexes,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("ошибка создания таблиц: %w", err)
		}
	}

	return nil
}

// CreateTables создает таблицы если они не существуют
const createProductTypesTable = `
CREATE TABLE IF NOT EXISTS product_types (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	coefficient DECIMAL(10,6) NOT NULL DEFAULT 1.0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);`

const createMaterialTypesTable = `
CREATE TABLE IF NOT EXISTS material_types (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	waste_percentage DECIMAL(5,2) NOT NULL DEFAULT 0.0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);`

const createMeasurementUnitsTable = `
CREATE TABLE IF NOT EXISTS measurement_units (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE,
	abbreviation VARCHAR(20) NOT NULL UNIQUE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);`

const createMaterialsTable = `
CREATE TABLE IF NOT EXISTS materials (
	id SERIAL PRIMARY KEY,
	article VARCHAR(100) NOT NULL UNIQUE,
	material_type_id INTEGER NOT NULL REFERENCES material_types(id),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	measurement_unit_id INTEGER NOT NULL REFERENCES measurement_units(id),
	package_quantity DECIMAL(10,3) NOT NULL DEFAULT 1.0,
	cost_per_unit DECIMAL(12,2) NOT NULL,
	stock_quantity DECIMAL(10,3) NOT NULL DEFAULT 0.0,
	min_stock_quantity DECIMAL(10,3) NOT NULL DEFAULT 0.0,
	image_path VARCHAR(500),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);`

const createProductsTable = `
CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,
	article VARCHAR(100) NOT NULL UNIQUE,
	product_type_id INTEGER NOT NULL REFERENCES product_types(id),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	image_path VARCHAR(500),
	min_partner_price DECIMAL(12,2) NOT NULL,
	package_dimensions JSONB,
	weights JSONB,
	quality_certificate_path VARCHAR(500),
	standard_number VARCHAR(100),
	production_time JSONB,
	cost_price DECIMAL(12,2),
	workshop_number VARCHAR(50),
	required_workers INTEGER,
	roll_width DECIMAL(8,3),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);`

const createProductMaterialsTable = `
CREATE TABLE IF NOT EXISTS product_materials (
	id SERIAL PRIMARY KEY,
	product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
	material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
	quantity_per_unit DECIMAL(10,6) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(product_id, material_id)
);`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_materials_article ON materials(article);
CREATE INDEX IF NOT EXISTS idx_materials_material_type ON materials(material_type_id);
CREATE INDEX IF NOT EXISTS idx_materials_stock ON materials(stock_quantity);
CREATE INDEX IF NOT EXISTS idx_products_article ON products(article);
CREATE INDEX IF NOT EXISTS idx_products_product_type ON products(product_type_id);
CREATE INDEX IF NOT EXISTS idx_product_materials_product ON product_materials(product_id);
CREATE INDEX IF NOT EXISTS idx_product_materials_material ON product_materials(material_id);`

// InsertTestData вставляет тестовые данные
func InsertTestData(db *sql.DB) error {
	// Вставляем типы продукции
	productTypes := []struct {
		name        string
		coefficient float64
	}{
		{"Обои виниловые", 1.05},
		{"Обои флизелиновые", 1.03},
		{"Обои бумажные", 1.02},
		{"Фотообои", 1.1},
	}

	for _, pt := range productTypes {
		_, err := db.Exec(`
			INSERT INTO product_types (name, coefficient) 
			VALUES ($1, $2) 
			ON CONFLICT (name) DO NOTHING`,
			pt.name, pt.coefficient)
		if err != nil {
			return fmt.Errorf("ошибка вставки типа продукции: %w", err)
		}
	}

	// Вставляем типы материалов
	materialTypes := []struct {
		name            string
		wastePercentage float64
	}{
		{"Винил", 5.0},
		{"Флизелин", 3.0},
		{"Бумага", 2.0},
		{"Краска", 1.0},
		{"Клей", 0.5},
	}

	for _, mt := range materialTypes {
		_, err := db.Exec(`
			INSERT INTO material_types (name, waste_percentage) 
			VALUES ($1, $2) 
			ON CONFLICT (name) DO NOTHING`,
			mt.name, mt.wastePercentage)
		if err != nil {
			return fmt.Errorf("ошибка вставки типа материала: %w", err)
		}
	}

	// Вставляем единицы измерения
	units := []struct {
		name         string
		abbreviation string
	}{
		{"Квадратный метр", "м²"},
		{"Погонный метр", "пог.м"},
		{"Литр", "л"},
		{"Килограмм", "кг"},
		{"Штука", "шт"},
		{"Рулон", "рул"},
	}

	for _, unit := range units {
		_, err := db.Exec(`
			INSERT INTO measurement_units (name, abbreviation) 
			VALUES ($1, $2) 
			ON CONFLICT (name) DO NOTHING`,
			unit.name, unit.abbreviation)
		if err != nil {
			return fmt.Errorf("ошибка вставки единицы измерения: %w", err)
		}
	}

	return nil
}

// RunMigrations запускает все миграции
func RunMigrations(db *sql.DB) error {
	if err := CreateTables(db); err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	if err := InsertTestData(db); err != nil {
		return fmt.Errorf("ошибка вставки тестовых данных: %w", err)
	}

	return nil
}

// HealthCheck проверяет состояние подключения к БД
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("база данных недоступна: %w", err)
	}

	return nil
}
