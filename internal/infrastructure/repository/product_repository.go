package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"

	"github.com/Masterminds/squirrel"
)

// ProductRepository реализует интерфейс repository.ProductRepository
type ProductRepository struct {
	*BaseRepository
}

// NewProductRepository создает новый репозиторий продукции
func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &ProductRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create создает новую продукцию
func (r *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	query := r.Builder().
		Insert(TableNames.Products).
		Columns(
			ProductColumns.Article,
			ProductColumns.ProductTypeID,
			ProductColumns.Name,
			ProductColumns.Description,
			ProductColumns.ImagePath,
			ProductColumns.MinPartnerPrice,
			ProductColumns.PackageDimensions,
			ProductColumns.Weights,
			ProductColumns.QualityCertificatePath,
			ProductColumns.StandardNumber,
			ProductColumns.ProductionTime,
			ProductColumns.CostPrice,
			ProductColumns.WorkshopNumber,
			ProductColumns.RequiredWorkers,
			ProductColumns.RollWidth,
			ProductColumns.CreatedAt,
			ProductColumns.UpdatedAt,
		).
		Values(
			product.Article(),
			int(product.ProductTypeID()),
			product.Name(),
			product.Description(),
			product.ImagePath(),
			float64(product.MinPartnerPrice()),
			r.serializePackageDimensions(product.PackageDimensions()),
			r.serializeWeights(product.Weights()),
			product.QualityCertificatePath(),
			product.StandardNumber(),
			r.serializeProductionTime(product.ProductionTime()),
			r.serializeCostPrice(product.CostPrice()),
			product.WorkshopNumber(),
			product.RequiredWorkers(),
			product.RollWidth(),
			time.Now(),
			time.Now(),
		).
		Suffix("RETURNING id")

	var id int
	err := r.QueryRow(query).Scan(&id)
	if err != nil {
		return fmt.Errorf("ошибка создания продукции: %w", err)
	}

	// Устанавливаем ID в entity
	// Примечание: в реальном проекте нужно добавить метод SetID в entity
	return nil
}

// GetByID возвращает продукцию по ID
func (r *ProductRepository) GetByID(ctx context.Context, id entity.ID) (*entity.Product, error) {
	query := r.Builder().
		Select(
			ProductColumns.ID,
			ProductColumns.Article,
			ProductColumns.ProductTypeID,
			ProductColumns.Name,
			ProductColumns.Description,
			ProductColumns.ImagePath,
			ProductColumns.MinPartnerPrice,
			ProductColumns.PackageDimensions,
			ProductColumns.Weights,
			ProductColumns.QualityCertificatePath,
			ProductColumns.StandardNumber,
			ProductColumns.ProductionTime,
			ProductColumns.CostPrice,
			ProductColumns.WorkshopNumber,
			ProductColumns.RequiredWorkers,
			ProductColumns.RollWidth,
			ProductColumns.CreatedAt,
			ProductColumns.UpdatedAt,
		).
		From(TableNames.Products).
		Where(squirrel.Eq{ProductColumns.ID: int(id)}).
		Limit(1)

	row := r.QueryRow(query)
	return r.scanProduct(row)
}

// GetAll возвращает все продукции с пагинацией
func (r *ProductRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	query := r.Builder().
		Select(
			ProductColumns.ID,
			ProductColumns.Article,
			ProductColumns.ProductTypeID,
			ProductColumns.Name,
			ProductColumns.Description,
			ProductColumns.ImagePath,
			ProductColumns.MinPartnerPrice,
			ProductColumns.PackageDimensions,
			ProductColumns.Weights,
			ProductColumns.QualityCertificatePath,
			ProductColumns.StandardNumber,
			ProductColumns.ProductionTime,
			ProductColumns.CostPrice,
			ProductColumns.WorkshopNumber,
			ProductColumns.RequiredWorkers,
			ProductColumns.RollWidth,
			ProductColumns.CreatedAt,
			ProductColumns.UpdatedAt,
		).
		From(TableNames.Products).
		OrderBy(ProductColumns.CreatedAt + " DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка продукции: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product, err := r.scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования продукции: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по результатам: %w", err)
	}

	return products, nil
}

// Update обновляет продукцию
func (r *ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	query := r.Builder().
		Update(TableNames.Products).
		Set(ProductColumns.Name, product.Name()).
		Set(ProductColumns.Description, product.Description()).
		Set(ProductColumns.ImagePath, product.ImagePath()).
		Set(ProductColumns.MinPartnerPrice, float64(product.MinPartnerPrice())).
		Set(ProductColumns.PackageDimensions, r.serializePackageDimensions(product.PackageDimensions())).
		Set(ProductColumns.Weights, r.serializeWeights(product.Weights())).
		Set(ProductColumns.QualityCertificatePath, product.QualityCertificatePath()).
		Set(ProductColumns.StandardNumber, product.StandardNumber()).
		Set(ProductColumns.ProductionTime, r.serializeProductionTime(product.ProductionTime())).
		Set(ProductColumns.CostPrice, r.serializeCostPrice(product.CostPrice())).
		Set(ProductColumns.WorkshopNumber, product.WorkshopNumber()).
		Set(ProductColumns.RequiredWorkers, product.RequiredWorkers()).
		Set(ProductColumns.RollWidth, product.RollWidth()).
		Set(ProductColumns.UpdatedAt, time.Now()).
		Where(squirrel.Eq{ProductColumns.ID: int(product.ID())})

	result, err := r.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("ошибка обновления продукции: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("продукция с ID %d не найдена", int(product.ID()))
	}

	return nil
}

// Delete удаляет продукцию
func (r *ProductRepository) Delete(ctx context.Context, id entity.ID) error {
	query := r.Builder().
		Delete(TableNames.Products).
		Where(squirrel.Eq{ProductColumns.ID: int(id)})

	result, err := r.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("ошибка удаления продукции: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("продукция с ID %d не найдена", int(id))
	}

	return nil
}

// GetByArticle возвращает продукцию по артикулу
func (r *ProductRepository) GetByArticle(ctx context.Context, article string) (*entity.Product, error) {
	query := r.Builder().
		Select(
			ProductColumns.ID,
			ProductColumns.Article,
			ProductColumns.ProductTypeID,
			ProductColumns.Name,
			ProductColumns.Description,
			ProductColumns.ImagePath,
			ProductColumns.MinPartnerPrice,
			ProductColumns.PackageDimensions,
			ProductColumns.Weights,
			ProductColumns.QualityCertificatePath,
			ProductColumns.StandardNumber,
			ProductColumns.ProductionTime,
			ProductColumns.CostPrice,
			ProductColumns.WorkshopNumber,
			ProductColumns.RequiredWorkers,
			ProductColumns.RollWidth,
			ProductColumns.CreatedAt,
			ProductColumns.UpdatedAt,
		).
		From(TableNames.Products).
		Where(squirrel.Eq{ProductColumns.Article: article}).
		Limit(1)

	row := r.QueryRow(query)
	return r.scanProduct(row)
}

// ExistsByArticle проверяет существование продукции по артикулу
func (r *ProductRepository) ExistsByArticle(ctx context.Context, article string) (bool, error) {
	query := r.Builder().
		Select("1").
		From(TableNames.Products).
		Where(squirrel.Eq{ProductColumns.Article: article}).
		Limit(1)

	var exists int
	err := r.QueryRow(query).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("ошибка проверки существования продукции: %w", err)
	}

	return true, nil
}

// Count возвращает общее количество продукций
func (r *ProductRepository) Count(ctx context.Context) (int, error) {
	query := r.BuildCountQuery(TableNames.Products)

	var count int
	err := r.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчета продукции: %w", err)
	}

	return count, nil
}

// Search ищет продукции по запросу
func (r *ProductRepository) Search(ctx context.Context, searchQuery string, limit, offset int) ([]*entity.Product, error) {
	query := r.Builder().
		Select(
			ProductColumns.ID,
			ProductColumns.Article,
			ProductColumns.ProductTypeID,
			ProductColumns.Name,
			ProductColumns.Description,
			ProductColumns.ImagePath,
			ProductColumns.MinPartnerPrice,
			ProductColumns.PackageDimensions,
			ProductColumns.Weights,
			ProductColumns.QualityCertificatePath,
			ProductColumns.StandardNumber,
			ProductColumns.ProductionTime,
			ProductColumns.CostPrice,
			ProductColumns.WorkshopNumber,
			ProductColumns.RequiredWorkers,
			ProductColumns.RollWidth,
			ProductColumns.CreatedAt,
			ProductColumns.UpdatedAt,
		).
		From(TableNames.Products).
		Where(squirrel.Or{
			squirrel.Like{ProductColumns.Name: "%" + searchQuery + "%"},
			squirrel.Like{ProductColumns.Article: "%" + searchQuery + "%"},
			squirrel.Like{ProductColumns.Description: "%" + searchQuery + "%"},
		}).
		OrderBy(ProductColumns.CreatedAt + " DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска продукции: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product, err := r.scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования продукции: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// GetMaterials возвращает материалы для продукции
func (r *ProductRepository) GetMaterials(ctx context.Context, productID entity.ID) ([]*entity.ProductMaterial, error) {
	query := r.Builder().
		Select(
			"pm."+ProductMaterialColumns.ID,
			"pm."+ProductMaterialColumns.ProductID,
			"pm."+ProductMaterialColumns.MaterialID,
			"pm."+ProductMaterialColumns.QuantityPerUnit,
			"pm."+ProductMaterialColumns.CreatedAt,
			"pm."+ProductMaterialColumns.UpdatedAt,
			// Материал
			"m."+MaterialColumns.ID,
			"m."+MaterialColumns.Article,
			"m."+MaterialColumns.Name,
			"m."+MaterialColumns.CostPerUnit,
		).
		From(TableNames.ProductMaterials + " pm").
		Join(TableNames.Materials + " m ON pm." + ProductMaterialColumns.MaterialID + " = m." + MaterialColumns.ID).
		Where(squirrel.Eq{"pm." + ProductMaterialColumns.ProductID: int(productID)}).
		OrderBy("pm." + ProductMaterialColumns.CreatedAt)

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов продукции: %w", err)
	}
	defer rows.Close()

	var productMaterials []*entity.ProductMaterial
	for rows.Next() {
		// Сканируем данные
		var pm struct {
			ID              int
			ProductID       int
			MaterialID      int
			QuantityPerUnit float64
			CreatedAt       time.Time
			UpdatedAt       time.Time
		}
		var material struct {
			ID          int
			Article     string
			Name        string
			CostPerUnit float64
		}

		err := rows.Scan(
			&pm.ID,
			&pm.ProductID,
			&pm.MaterialID,
			&pm.QuantityPerUnit,
			&pm.CreatedAt,
			&pm.UpdatedAt,
			&material.ID,
			&material.Article,
			&material.Name,
			&material.CostPerUnit,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала продукции: %w", err)
		}

		// Создаем entity (упрощенная версия для демонстрации)
		// В реальном проекте нужно создать полноценную entity
		// productMaterial := entity.NewProductMaterial(...)
		// productMaterials = append(productMaterials, productMaterial)
	}

	return productMaterials, nil
}

// scanProduct сканирует данные продукции из Row или Rows
func (r *ProductRepository) scanProduct(scanner interface{}) (*entity.Product, error) {
	var (
		id                       int
		article                  string
		productTypeID            int
		name                     string
		description              sql.NullString
		imagePath                sql.NullString
		minPartnerPrice          float64
		packageDimensions        sql.NullString
		weights                  sql.NullString
		qualityCertificatePath   sql.NullString
		standardNumber           sql.NullString
		productionTime           sql.NullString
		costPrice                sql.NullFloat64
		workshopNumber           sql.NullString
		requiredWorkers          sql.NullInt32
		rollWidth                sql.NullFloat64
		createdAt                time.Time
		updatedAt                time.Time
	)

	var err error
	switch s := scanner.(type) {
	case *sql.Row:
		err = s.Scan(
			&id, &article, &productTypeID, &name, &description, &imagePath,
			&minPartnerPrice, &packageDimensions, &weights, &qualityCertificatePath,
			&standardNumber, &productionTime, &costPrice, &workshopNumber,
			&requiredWorkers, &rollWidth, &createdAt, &updatedAt,
		)
	case *sql.Rows:
		err = s.Scan(
			&id, &article, &productTypeID, &name, &description, &imagePath,
			&minPartnerPrice, &packageDimensions, &weights, &qualityCertificatePath,
			&standardNumber, &productionTime, &costPrice, &workshopNumber,
			&requiredWorkers, &rollWidth, &createdAt, &updatedAt,
		)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип scanner")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("продукция не найдена")
		}
		return nil, fmt.Errorf("ошибка сканирования продукции: %w", err)
	}

	// Создаем продукцию через конструктор entity
	product, err := entity.NewProduct(
		article,
		entity.ID(productTypeID),
		name,
		entity.Money(minPartnerPrice),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания entity продукции: %w", err)
	}

	// Устанавливаем дополнительные поля, если они есть
	if description.Valid {
		// product.SetDescription(&description.String)
	}

	// Примечание: В реальном проекте нужно добавить методы для установки всех полей

	return product, nil
}

// Вспомогательные методы для сериализации сложных типов в JSON

func (r *ProductRepository) serializePackageDimensions(dimensions *entity.PackageDimensions) interface{} {
	if dimensions == nil {
		return nil
	}
	data, _ := json.Marshal(dimensions)
	return string(data)
}

func (r *ProductRepository) serializeWeights(weights *entity.ProductWeights) interface{} {
	if weights == nil {
		return nil
	}
	data, _ := json.Marshal(weights)
	return string(data)
}

func (r *ProductRepository) serializeProductionTime(time *entity.ProductionTime) interface{} {
	if time == nil {
		return nil
	}
	data, _ := json.Marshal(time)
	return string(data)
}

func (r *ProductRepository) serializeCostPrice(price *entity.Money) interface{} {
	if price == nil {
		return nil
	}
	return float64(*price)
} 