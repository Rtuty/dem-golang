package repositories

import (
	"database/sql"
	"fmt"
	"strconv"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/repositories"
)

// materialRepositoryImpl реализует интерфейс MaterialRepository
type materialRepositoryImpl struct {
	db *sql.DB
}

// NewMaterialRepository создает новую реализацию репозитория материалов
func NewMaterialRepository(db *sql.DB) repositories.MaterialRepository {
	return &materialRepositoryImpl{db: db}
}

// GetAll возвращает список всех материалов
func (r *materialRepositoryImpl) GetAll() ([]entities.Material, error) {
	query := `
		SELECT 
			m.id, m.article, m.material_type_id, m.name, m.description,
			m.measurement_unit_id, m.package_quantity, m.cost_per_unit,
			m.stock_quantity, m.min_stock_quantity, m.image_path,
			m.created_at, m.updated_at,
			mt.name as type_name, mt.waste_percentage,
			mu.name as unit_name, mu.abbreviation
		FROM materials m
		JOIN material_types mt ON m.material_type_id = mt.id
		JOIN measurement_units mu ON m.measurement_unit_id = mu.id
		ORDER BY m.name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса материалов: %w", err)
	}
	defer rows.Close()

	var materials []entities.Material
	for rows.Next() {
		var material entities.Material
		var typeName string
		var wastePercentage float64
		var unitName, unitAbbr string

		err := rows.Scan(
			&material.ID, &material.Article, &material.MaterialTypeID, &material.Name,
			&material.Description, &material.MeasurementUnitID, &material.PackageQuantity,
			&material.CostPerUnit, &material.StockQuantity, &material.MinStockQuantity,
			&material.ImagePath, &material.CreatedAt, &material.UpdatedAt,
			&typeName, &wastePercentage, &unitName, &unitAbbr,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}

		// Заполняем связанные данные
		material.MaterialType = &entities.MaterialType{
			ID:              material.MaterialTypeID,
			Name:            typeName,
			WastePercentage: wastePercentage,
		}

		material.MeasurementUnit = &entities.MeasurementUnit{
			ID:           material.MeasurementUnitID,
			Name:         unitName,
			Abbreviation: unitAbbr,
		}

		materials = append(materials, material)
	}

	return materials, nil
}

// GetByID возвращает материал по ID
func (r *materialRepositoryImpl) GetByID(id int) (*entities.Material, error) {
	query := `
		SELECT 
			m.id, m.article, m.material_type_id, m.name, m.description,
			m.measurement_unit_id, m.package_quantity, m.cost_per_unit,
			m.stock_quantity, m.min_stock_quantity, m.image_path,
			m.created_at, m.updated_at,
			mt.name as type_name, mt.waste_percentage,
			mu.name as unit_name, mu.abbreviation
		FROM materials m
		JOIN material_types mt ON m.material_type_id = mt.id
		JOIN measurement_units mu ON m.measurement_unit_id = mu.id
		WHERE m.id = $1
	`

	var material entities.Material
	var typeName string
	var wastePercentage float64
	var unitName, unitAbbr string

	err := r.db.QueryRow(query, id).Scan(
		&material.ID, &material.Article, &material.MaterialTypeID, &material.Name,
		&material.Description, &material.MeasurementUnitID, &material.PackageQuantity,
		&material.CostPerUnit, &material.StockQuantity, &material.MinStockQuantity,
		&material.ImagePath, &material.CreatedAt, &material.UpdatedAt,
		&typeName, &wastePercentage, &unitName, &unitAbbr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.NewNotFoundError("материал", strconv.Itoa(id))
		}
		return nil, fmt.Errorf("ошибка получения материала: %w", err)
	}

	// Заполняем связанные данные
	material.MaterialType = &entities.MaterialType{
		ID:              material.MaterialTypeID,
		Name:            typeName,
		WastePercentage: wastePercentage,
	}

	material.MeasurementUnit = &entities.MeasurementUnit{
		ID:           material.MeasurementUnitID,
		Name:         unitName,
		Abbreviation: unitAbbr,
	}

	return &material, nil
}

// GetMaterialTypeByID возвращает тип материала по ID
func (r *materialRepositoryImpl) GetMaterialTypeByID(id int) (*entities.MaterialType, error) {
	query := `
		SELECT id, name, waste_percentage, created_at, updated_at 
		FROM material_types 
		WHERE id = $1
	`

	var materialType entities.MaterialType
	err := r.db.QueryRow(query, id).Scan(
		&materialType.ID, &materialType.Name, &materialType.WastePercentage,
		&materialType.CreatedAt, &materialType.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.NewNotFoundError("тип материала", strconv.Itoa(id))
		}
		return nil, fmt.Errorf("ошибка получения типа материала: %w", err)
	}

	return &materialType, nil
}

// GetProductTypeByID возвращает тип продукции по ID
func (r *materialRepositoryImpl) GetProductTypeByID(id int) (*entities.ProductType, error) {
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

// GetMaterialTypes возвращает все типы материалов
func (r *materialRepositoryImpl) GetMaterialTypes() ([]entities.MaterialType, error) {
	query := "SELECT id, name, waste_percentage, created_at, updated_at FROM material_types ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса типов материалов: %w", err)
	}
	defer rows.Close()

	var types []entities.MaterialType
	for rows.Next() {
		var materialType entities.MaterialType
		err := rows.Scan(&materialType.ID, &materialType.Name, &materialType.WastePercentage,
			&materialType.CreatedAt, &materialType.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования типа материала: %w", err)
		}
		types = append(types, materialType)
	}

	return types, nil
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (r *materialRepositoryImpl) GetMaterialsForProduct(productID int) ([]entities.Material, error) {
	query := `
		SELECT 
			m.id, m.article, m.material_type_id, m.name, m.description,
			m.measurement_unit_id, m.package_quantity, m.cost_per_unit,
			m.stock_quantity, m.min_stock_quantity, m.image_path,
			m.created_at, m.updated_at,
			mt.name as type_name, mt.waste_percentage,
			mu.name as unit_name, mu.abbreviation
		FROM product_materials pm
		JOIN materials m ON pm.material_id = m.id
		JOIN material_types mt ON m.material_type_id = mt.id
		JOIN measurement_units mu ON m.measurement_unit_id = mu.id
		WHERE pm.product_id = $1
		ORDER BY m.name
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса материалов: %w", err)
	}
	defer rows.Close()

	var materials []entities.Material
	for rows.Next() {
		var material entities.Material
		var typeName string
		var wastePercentage float64
		var unitName, unitAbbr string

		err := rows.Scan(
			&material.ID, &material.Article, &material.MaterialTypeID, &material.Name,
			&material.Description, &material.MeasurementUnitID, &material.PackageQuantity,
			&material.CostPerUnit, &material.StockQuantity, &material.MinStockQuantity,
			&material.ImagePath, &material.CreatedAt, &material.UpdatedAt,
			&typeName, &wastePercentage, &unitName, &unitAbbr,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}

		// Заполняем связанные данные
		material.MaterialType = &entities.MaterialType{
			ID:              material.MaterialTypeID,
			Name:            typeName,
			WastePercentage: wastePercentage,
		}

		material.MeasurementUnit = &entities.MeasurementUnit{
			ID:           material.MeasurementUnitID,
			Name:         unitName,
			Abbreviation: unitAbbr,
		}

		materials = append(materials, material)
	}

	return materials, nil
}
