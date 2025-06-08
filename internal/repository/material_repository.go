package repository

import (
	"database/sql"
	"fmt"

	"wallpaper-system/internal/models"
)

// MaterialRepository представляет репозиторий для работы с материалами
type MaterialRepository struct {
	db *sql.DB
}

// NewMaterialRepository создает новый репозиторий материалов
func NewMaterialRepository(db *sql.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции с детальной информацией
func (r *MaterialRepository) GetMaterialsForProduct(productID int) ([]models.MaterialForProduct, error) {
	query := `
		SELECT 
			m.id, m.name, m.article, pm.quantity_per_unit, m.cost_per_unit,
			mu.abbreviation
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

	var materials []models.MaterialForProduct
	for rows.Next() {
		var material models.MaterialForProduct

		err := rows.Scan(
			&material.ID, &material.Name, &material.Article,
			&material.QuantityPerUnit, &material.CostPerUnit, &material.UnitAbbreviation,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}

		// Рассчитываем общую стоимость
		material.TotalCost = material.QuantityPerUnit * material.CostPerUnit

		materials = append(materials, material)
	}

	return materials, nil
}

// GetMaterialTypeByID возвращает тип материала по ID
func (r *MaterialRepository) GetMaterialTypeByID(id int) (*models.MaterialType, error) {
	query := `
		SELECT id, name, waste_percentage, created_at, updated_at 
		FROM material_types 
		WHERE id = $1
	`

	var materialType models.MaterialType
	err := r.db.QueryRow(query, id).Scan(
		&materialType.ID, &materialType.Name, &materialType.WastePercentage,
		&materialType.CreatedAt, &materialType.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("тип материала с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения типа материала: %w", err)
	}

	return &materialType, nil
}

// GetProductTypeByID возвращает тип продукции по ID
func (r *MaterialRepository) GetProductTypeByID(id int) (*models.ProductType, error) {
	query := `
		SELECT id, name, coefficient, created_at, updated_at 
		FROM product_types 
		WHERE id = $1
	`

	var productType models.ProductType
	err := r.db.QueryRow(query, id).Scan(
		&productType.ID, &productType.Name, &productType.Coefficient,
		&productType.CreatedAt, &productType.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("тип продукции с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения типа продукции: %w", err)
	}

	return &productType, nil
}

// GetAllMaterials возвращает все материалы
func (r *MaterialRepository) GetAllMaterials() ([]models.Material, error) {
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

	var materials []models.Material
	for rows.Next() {
		var material models.Material
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
		material.MaterialType = &models.MaterialType{
			ID:              material.MaterialTypeID,
			Name:            typeName,
			WastePercentage: wastePercentage,
		}

		material.MeasurementUnit = &models.MeasurementUnit{
			ID:           material.MeasurementUnitID,
			Name:         unitName,
			Abbreviation: unitAbbr,
		}

		materials = append(materials, material)
	}

	return materials, nil
}

// GetMaterialTypes возвращает все типы материалов
func (r *MaterialRepository) GetMaterialTypes() ([]models.MaterialType, error) {
	query := "SELECT id, name, waste_percentage, created_at, updated_at FROM material_types ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса типов материалов: %w", err)
	}
	defer rows.Close()

	var types []models.MaterialType
	for rows.Next() {
		var materialType models.MaterialType
		err := rows.Scan(&materialType.ID, &materialType.Name, &materialType.WastePercentage,
			&materialType.CreatedAt, &materialType.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования типа материала: %w", err)
		}
		types = append(types, materialType)
	}

	return types, nil
}
