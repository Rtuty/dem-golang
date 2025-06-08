package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"wallpaper-system/internal/models"
)

// ProductRepository представляет репозиторий для работы с продукцией
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository создает новый репозиторий продукции
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll возвращает список всей продукции с типами
func (r *ProductRepository) GetAll() ([]models.ProductListItem, error) {
	query := `
		SELECT 
			p.id, p.article, pt.name as type_name, p.name,
			p.min_partner_price, p.roll_width
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var products []models.ProductListItem
	for rows.Next() {
		var product models.ProductListItem
		err := rows.Scan(
			&product.ID, &product.Article, &product.TypeName, &product.Name,
			&product.MinPartnerPrice, &product.RollWidth,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// GetByID возвращает продукцию по ID с материалами
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
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

	var product models.Product
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
			return nil, fmt.Errorf("продукция с ID %d не найдена", id)
		}
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Заполняем информацию о типе
	product.ProductType = &models.ProductType{
		ID:          product.ProductTypeID,
		Name:        typeName,
		Coefficient: typeCoefficient,
	}

	// Получаем материалы для продукции
	materials, err := r.GetMaterialsForProduct(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов: %w", err)
	}
	product.Materials = materials

	return &product, nil
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (r *ProductRepository) GetMaterialsForProduct(productID int) ([]models.ProductMaterial, error) {
	query := `
		SELECT 
			pm.id, pm.product_id, pm.material_id, pm.quantity_per_unit, pm.created_at,
			m.article, m.name, m.cost_per_unit,
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

	var materials []models.ProductMaterial
	for rows.Next() {
		var pm models.ProductMaterial
		var material models.Material
		var unitAbbr string

		err := rows.Scan(
			&pm.ID, &pm.ProductID, &pm.MaterialID, &pm.QuantityPerUnit, &pm.CreatedAt,
			&material.Article, &material.Name, &material.CostPerUnit, &unitAbbr,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}

		material.ID = pm.MaterialID
		material.MeasurementUnit = &models.MeasurementUnit{Abbreviation: unitAbbr}
		pm.Material = &material

		materials = append(materials, pm)
	}

	return materials, nil
}

// Create создает новую продукцию
func (r *ProductRepository) Create(req *models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products (article, product_type_id, name, description, min_partner_price, roll_width)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var product models.Product
	err := r.db.QueryRow(query, req.Article, req.ProductTypeID, req.Name,
		req.Description, req.MinPartnerPrice, req.RollWidth).Scan(
		&product.ID, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания продукции: %w", err)
	}

	// Заполняем остальные поля
	product.Article = req.Article
	product.ProductTypeID = req.ProductTypeID
	product.Name = req.Name
	product.Description = req.Description
	product.MinPartnerPrice = req.MinPartnerPrice
	product.RollWidth = req.RollWidth

	return &product, nil
}

// Update обновляет продукцию
func (r *ProductRepository) Update(id int, req *models.UpdateProductRequest) error {
	// Построение динамического запроса обновления
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Article != nil {
		setParts = append(setParts, fmt.Sprintf("article = $%d", argIndex))
		args = append(args, *req.Article)
		argIndex++
	}
	if req.ProductTypeID != nil {
		setParts = append(setParts, fmt.Sprintf("product_type_id = $%d", argIndex))
		args = append(args, *req.ProductTypeID)
		argIndex++
	}
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.MinPartnerPrice != nil {
		setParts = append(setParts, fmt.Sprintf("min_partner_price = $%d", argIndex))
		args = append(args, *req.MinPartnerPrice)
		argIndex++
	}
	if req.RollWidth != nil {
		setParts = append(setParts, fmt.Sprintf("roll_width = $%d", argIndex))
		args = append(args, *req.RollWidth)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("нет полей для обновления")
	}

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления продукции: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("продукция с ID %d не найдена", id)
	}

	return nil
}

// Delete удаляет продукцию
func (r *ProductRepository) Delete(id int) error {
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
		return fmt.Errorf("продукция с ID %d не найдена", id)
	}

	return nil
}

// GetProductTypes возвращает все типы продукции
func (r *ProductRepository) GetProductTypes() ([]models.ProductType, error) {
	query := "SELECT id, name, coefficient, created_at, updated_at FROM product_types ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса типов продукции: %w", err)
	}
	defer rows.Close()

	var types []models.ProductType
	for rows.Next() {
		var productType models.ProductType
		err := rows.Scan(&productType.ID, &productType.Name, &productType.Coefficient,
			&productType.CreatedAt, &productType.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования типа продукции: %w", err)
		}
		types = append(types, productType)
	}

	return types, nil
}
