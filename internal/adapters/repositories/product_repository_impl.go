package repositories

import (
	"database/sql"
	"fmt"
	"strconv"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/repositories"
)

// productRepositoryImpl реализует интерфейс ProductRepository
type productRepositoryImpl struct {
	db *sql.DB
}

// NewProductRepository создает новую реализацию репозитория продукции
func NewProductRepository(db *sql.DB) repositories.ProductRepository {
	return &productRepositoryImpl{db: db}
}

// GetAll возвращает список всей продукции
func (r *productRepositoryImpl) GetAll() ([]entities.Product, error) {
	query := `
		SELECT 
			p.id, p.article, p.product_type_id, p.name, p.description,
			p.image_path, p.min_partner_price, p.package_length, p.package_width,
			p.package_height, p.weight_without_package, p.weight_with_package,
			p.quality_certificate_path, p.standard_number, p.production_time_hours,
			p.cost_price, p.workshop_number, p.required_workers, p.roll_width,
			p.created_at, p.updated_at,
			pt.name as type_name, pt.coefficient as type_coefficient
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		var typeName string
		var typeCoefficient float64

		err := rows.Scan(
			&product.ID, &product.Article, &product.ProductTypeID, &product.Name,
			&product.Description, &product.ImagePath, &product.MinPartnerPrice,
			&product.PackageLength, &product.PackageWidth, &product.PackageHeight,
			&product.WeightWithoutPackage, &product.WeightWithPackage,
			&product.QualityCertificatePath, &product.StandardNumber,
			&product.ProductionTimeHours, &product.CostPrice, &product.WorkshopNumber,
			&product.RequiredWorkers, &product.RollWidth, &product.CreatedAt,
			&product.UpdatedAt, &typeName, &typeCoefficient,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		// Заполняем тип продукции
		product.ProductType = &entities.ProductType{
			ID:          product.ProductTypeID,
			Name:        typeName,
			Coefficient: typeCoefficient,
		}

		products = append(products, product)
	}

	return products, nil
}

// GetByID возвращает продукцию по ID
func (r *productRepositoryImpl) GetByID(id int) (*entities.Product, error) {
	query := `
		SELECT 
			p.id, p.article, p.product_type_id, p.name, p.description,
			p.image_path, p.min_partner_price, p.package_length, p.package_width,
			p.package_height, p.weight_without_package, p.weight_with_package,
			p.quality_certificate_path, p.standard_number, p.production_time_hours,
			p.cost_price, p.workshop_number, p.required_workers, p.roll_width,
			p.created_at, p.updated_at,
			pt.name as type_name, pt.coefficient as type_coefficient
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		WHERE p.id = $1
	`

	var product entities.Product
	var typeName string
	var typeCoefficient float64

	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Article, &product.ProductTypeID, &product.Name,
		&product.Description, &product.ImagePath, &product.MinPartnerPrice,
		&product.PackageLength, &product.PackageWidth, &product.PackageHeight,
		&product.WeightWithoutPackage, &product.WeightWithPackage,
		&product.QualityCertificatePath, &product.StandardNumber,
		&product.ProductionTimeHours, &product.CostPrice, &product.WorkshopNumber,
		&product.RequiredWorkers, &product.RollWidth, &product.CreatedAt,
		&product.UpdatedAt, &typeName, &typeCoefficient,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.NewNotFoundError("продукция", strconv.Itoa(id))
		}
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Заполняем тип продукции
	product.ProductType = &entities.ProductType{
		ID:          product.ProductTypeID,
		Name:        typeName,
		Coefficient: typeCoefficient,
	}

	// Получаем материалы
	materials, err := r.GetMaterialsForProduct(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов: %w", err)
	}
	product.Materials = materials

	return &product, nil
}

// Create создает новую продукцию
func (r *productRepositoryImpl) Create(product *entities.Product) error {
	query := `
		INSERT INTO products (article, product_type_id, name, description, min_partner_price, roll_width)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, product.Article, product.ProductTypeID, product.Name,
		product.Description, product.MinPartnerPrice, product.RollWidth).Scan(
		&product.ID, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("ошибка создания продукции: %w", err)
	}

	return nil
}

// Update обновляет продукцию
func (r *productRepositoryImpl) Update(product *entities.Product) error {
	query := `
		UPDATE products 
		SET article = $1, product_type_id = $2, name = $3, description = $4, 
		    min_partner_price = $5, roll_width = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	result, err := r.db.Exec(query, product.Article, product.ProductTypeID,
		product.Name, product.Description, product.MinPartnerPrice,
		product.RollWidth, product.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления продукции: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return entities.NewNotFoundError("продукция", strconv.Itoa(product.ID))
	}

	return nil
}

// Delete удаляет продукцию
func (r *productRepositoryImpl) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления продукции: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return entities.NewNotFoundError("продукция", strconv.Itoa(id))
	}

	return nil
}

// GetProductTypes возвращает все типы продукции
func (r *productRepositoryImpl) GetProductTypes() ([]entities.ProductType, error) {
	query := "SELECT id, name, coefficient, created_at, updated_at FROM product_types ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса типов продукции: %w", err)
	}
	defer rows.Close()

	var types []entities.ProductType
	for rows.Next() {
		var productType entities.ProductType
		err := rows.Scan(&productType.ID, &productType.Name, &productType.Coefficient,
			&productType.CreatedAt, &productType.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования типа продукции: %w", err)
		}
		types = append(types, productType)
	}

	return types, nil
}

// GetProductTypeByID возвращает тип продукции по ID
func (r *productRepositoryImpl) GetProductTypeByID(id int) (*entities.ProductType, error) {
	query := `
		SELECT id, name, coefficient, created_at, updated_at 
		FROM product_types 
		WHERE id = $1
	`

	var productType entities.ProductType
	err := r.db.QueryRow(query, id).Scan(
		&productType.ID, &productType.Name, &productType.Coefficient,
		&productType.CreatedAt, &productType.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.NewNotFoundError("тип продукции", strconv.Itoa(id))
		}
		return nil, fmt.Errorf("ошибка получения типа продукции: %w", err)
	}

	return &productType, nil
}

// GetMaterialsForProduct возвращает материалы для продукции
func (r *productRepositoryImpl) GetMaterialsForProduct(productID int) ([]entities.ProductMaterial, error) {
	query := `
		SELECT 
			pm.id, pm.product_id, pm.material_id, pm.quantity_per_unit, pm.created_at,
			m.id, m.article, m.material_type_id, m.name, m.description,
			m.measurement_unit_id, m.package_quantity, m.cost_per_unit,
			m.stock_quantity, m.min_stock_quantity, m.image_path,
			m.created_at, m.updated_at,
			mu.id, mu.name, mu.symbol as abbreviation, mu.created_at
		FROM product_materials pm
		JOIN materials m ON pm.material_id = m.id
		JOIN measurement_units mu ON m.measurement_unit_id = mu.id
		WHERE pm.product_id = $1
		ORDER BY m.name
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса материалов: %w", err)
	}
	defer rows.Close()

	var materials []entities.ProductMaterial
	for rows.Next() {
		var pm entities.ProductMaterial
		var material entities.Material
		var unit entities.MeasurementUnit

		err := rows.Scan(
			&pm.ID, &pm.ProductID, &pm.MaterialID, &pm.QuantityPerUnit, &pm.CreatedAt,
			&material.ID, &material.Article, &material.MaterialTypeID, &material.Name,
			&material.Description, &material.MeasurementUnitID, &material.PackageQuantity,
			&material.CostPerUnit, &material.StockQuantity, &material.MinStockQuantity,
			&material.ImagePath, &material.CreatedAt, &material.UpdatedAt,
			&unit.ID, &unit.Name, &unit.Abbreviation, &unit.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}

		material.MeasurementUnit = &unit
		pm.Material = &material

		materials = append(materials, pm)
	}

	return materials, nil
}
