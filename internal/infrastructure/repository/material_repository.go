package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"

	"github.com/Masterminds/squirrel"
)

// MaterialRepository реализует интерфейс repository.MaterialRepository
type MaterialRepository struct {
	*BaseRepository
}

// NewMaterialRepository создает новый репозиторий материалов
func NewMaterialRepository(db *sql.DB) repository.MaterialRepository {
	return &MaterialRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create создает новый материал
func (r *MaterialRepository) Create(ctx context.Context, material *entity.Material) error {
	query := r.Builder().
		Insert(TableNames.Materials).
		Columns(
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		Values(
			material.Article(),
			int(material.MaterialTypeID()),
			material.Name(),
			material.Description(),
			int(material.MeasurementUnitID()),
			material.PackageQuantity(),
			float64(material.CostPerUnit()),
			material.StockQuantity(),
			material.MinStockQuantity(),
			material.ImagePath(),
			time.Now(),
			time.Now(),
		).
		Suffix("RETURNING id")

	var id int
	err := r.QueryRow(query).Scan(&id)
	if err != nil {
		return fmt.Errorf("ошибка создания материала: %w", err)
	}

	return nil
}

// GetByID возвращает материал по ID
func (r *MaterialRepository) GetByID(ctx context.Context, id entity.ID) (*entity.Material, error) {
	query := r.Builder().
		Select(
			MaterialColumns.ID,
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		From(TableNames.Materials).
		Where(squirrel.Eq{MaterialColumns.ID: int(id)}).
		Limit(1)

	row := r.QueryRow(query)
	return r.scanMaterial(row)
}

// GetAll возвращает все материалы с пагинацией
func (r *MaterialRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Material, error) {
	query := r.Builder().
		Select(
			MaterialColumns.ID,
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		From(TableNames.Materials).
		OrderBy(MaterialColumns.CreatedAt + " DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка материалов: %w", err)
	}
	defer rows.Close()

	var materials []*entity.Material
	for rows.Next() {
		material, err := r.scanMaterial(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}
		materials = append(materials, material)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по результатам: %w", err)
	}

	return materials, nil
}

// Update обновляет материал
func (r *MaterialRepository) Update(ctx context.Context, material *entity.Material) error {
	query := r.Builder().
		Update(TableNames.Materials).
		Set(MaterialColumns.Name, material.Name()).
		Set(MaterialColumns.Description, material.Description()).
		Set(MaterialColumns.PackageQuantity, material.PackageQuantity()).
		Set(MaterialColumns.CostPerUnit, float64(material.CostPerUnit())).
		Set(MaterialColumns.StockQuantity, material.StockQuantity()).
		Set(MaterialColumns.MinStockQuantity, material.MinStockQuantity()).
		Set(MaterialColumns.ImagePath, material.ImagePath()).
		Set(MaterialColumns.UpdatedAt, time.Now()).
		Where(squirrel.Eq{MaterialColumns.ID: int(material.ID())})

	result, err := r.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("ошибка обновления материала: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("материал с ID %d не найден", int(material.ID()))
	}

	return nil
}

// Delete удаляет материал
func (r *MaterialRepository) Delete(ctx context.Context, id entity.ID) error {
	query := r.Builder().
		Delete(TableNames.Materials).
		Where(squirrel.Eq{MaterialColumns.ID: int(id)})

	result, err := r.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("ошибка удаления материала: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("материал с ID %d не найден", int(id))
	}

	return nil
}

// GetByArticle возвращает материал по артикулу
func (r *MaterialRepository) GetByArticle(ctx context.Context, article string) (*entity.Material, error) {
	query := r.Builder().
		Select(
			MaterialColumns.ID,
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		From(TableNames.Materials).
		Where(squirrel.Eq{MaterialColumns.Article: article}).
		Limit(1)

	row := r.QueryRow(query)
	return r.scanMaterial(row)
}

// ExistsByArticle проверяет существование материала по артикулу
func (r *MaterialRepository) ExistsByArticle(ctx context.Context, article string) (bool, error) {
	query := r.Builder().
		Select("1").
		From(TableNames.Materials).
		Where(squirrel.Eq{MaterialColumns.Article: article}).
		Limit(1)

	var exists int
	err := r.QueryRow(query).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("ошибка проверки существования материала: %w", err)
	}

	return true, nil
}

// Count возвращает общее количество материалов
func (r *MaterialRepository) Count(ctx context.Context) (int, error) {
	query := r.BuildCountQuery(TableNames.Materials)

	var count int
	err := r.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчета материалов: %w", err)
	}

	return count, nil
}

// Search ищет материалы по запросу
func (r *MaterialRepository) Search(ctx context.Context, searchQuery string, limit, offset int) ([]*entity.Material, error) {
	query := r.Builder().
		Select(
			MaterialColumns.ID,
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		From(TableNames.Materials).
		Where(squirrel.Or{
			squirrel.Like{MaterialColumns.Name: "%" + searchQuery + "%"},
			squirrel.Like{MaterialColumns.Article: "%" + searchQuery + "%"},
			squirrel.Like{MaterialColumns.Description: "%" + searchQuery + "%"},
		}).
		OrderBy(MaterialColumns.CreatedAt + " DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска материалов: %w", err)
	}
	defer rows.Close()

	var materials []*entity.Material
	for rows.Next() {
		material, err := r.scanMaterial(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}
		materials = append(materials, material)
	}

	return materials, nil
}

// GetLowStockMaterials возвращает материалы с низким остатком
func (r *MaterialRepository) GetLowStockMaterials(ctx context.Context) ([]*entity.Material, error) {
	query := r.Builder().
		Select(
			MaterialColumns.ID,
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		From(TableNames.Materials).
		Where(squirrel.LtOrEq{MaterialColumns.StockQuantity: squirrel.Expr(MaterialColumns.MinStockQuantity)}).
		OrderBy(MaterialColumns.StockQuantity + " ASC")

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов с низким остатком: %w", err)
	}
	defer rows.Close()

	var materials []*entity.Material
	for rows.Next() {
		material, err := r.scanMaterial(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}
		materials = append(materials, material)
	}

	return materials, nil
}

// GetByMaterialType возвращает материалы по типу
func (r *MaterialRepository) GetByMaterialType(ctx context.Context, materialTypeID entity.ID, limit, offset int) ([]*entity.Material, error) {
	query := r.Builder().
		Select(
			MaterialColumns.ID,
			MaterialColumns.Article,
			MaterialColumns.MaterialTypeID,
			MaterialColumns.Name,
			MaterialColumns.Description,
			MaterialColumns.MeasurementUnitID,
			MaterialColumns.PackageQuantity,
			MaterialColumns.CostPerUnit,
			MaterialColumns.StockQuantity,
			MaterialColumns.MinStockQuantity,
			MaterialColumns.ImagePath,
			MaterialColumns.CreatedAt,
			MaterialColumns.UpdatedAt,
		).
		From(TableNames.Materials).
		Where(squirrel.Eq{MaterialColumns.MaterialTypeID: int(materialTypeID)}).
		OrderBy(MaterialColumns.Name + " ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	rows, err := r.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов по типу: %w", err)
	}
	defer rows.Close()

	var materials []*entity.Material
	for rows.Next() {
		material, err := r.scanMaterial(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
		}
		materials = append(materials, material)
	}

	return materials, nil
}

// UpdateStock обновляет остаток материала
func (r *MaterialRepository) UpdateStock(ctx context.Context, id entity.ID, newQuantity float64) error {
	query := r.Builder().
		Update(TableNames.Materials).
		Set(MaterialColumns.StockQuantity, newQuantity).
		Set(MaterialColumns.UpdatedAt, time.Now()).
		Where(squirrel.Eq{MaterialColumns.ID: int(id)})

	result, err := r.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("ошибка обновления остатка материала: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("материал с ID %d не найден", int(id))
	}

	return nil
}

// AddStock увеличивает остаток материала
func (r *MaterialRepository) AddStock(ctx context.Context, id entity.ID, quantity float64) error {
	query := r.Builder().
		Update(TableNames.Materials).
		Set(MaterialColumns.StockQuantity, squirrel.Expr(MaterialColumns.StockQuantity+" + ?", quantity)).
		Set(MaterialColumns.UpdatedAt, time.Now()).
		Where(squirrel.Eq{MaterialColumns.ID: int(id)})

	result, err := r.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("ошибка увеличения остатка материала: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("материал с ID %d не найден", int(id))
	}

	return nil
}

// RemoveStock уменьшает остаток материала
func (r *MaterialRepository) RemoveStock(ctx context.Context, id entity.ID, quantity float64) error {
	// Сначала проверяем, что остатка достаточно
	checkQuery := r.Builder().
		Select(MaterialColumns.StockQuantity).
		From(TableNames.Materials).
		Where(squirrel.Eq{MaterialColumns.ID: int(id)})

	var currentStock float64
	err := r.QueryRow(checkQuery).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("материал с ID %d не найден", int(id))
		}
		return fmt.Errorf("ошибка проверки остатка материала: %w", err)
	}

	if currentStock < quantity {
		return fmt.Errorf("недостаточно материала на складе (есть: %.2f, требуется: %.2f)", currentStock, quantity)
	}

	// Уменьшаем остаток
	updateQuery := r.Builder().
		Update(TableNames.Materials).
		Set(MaterialColumns.StockQuantity, squirrel.Expr(MaterialColumns.StockQuantity+" - ?", quantity)).
		Set(MaterialColumns.UpdatedAt, time.Now()).
		Where(squirrel.Eq{MaterialColumns.ID: int(id)})

	result, err := r.ExecQuery(updateQuery)
	if err != nil {
		return fmt.Errorf("ошибка уменьшения остатка материала: %w", err)
	}

	rowsAffected, err := r.GetRowsAffected(result)
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("материал с ID %d не найден", int(id))
	}

	return nil
}

// scanMaterial сканирует данные материала из Row или Rows
func (r *MaterialRepository) scanMaterial(scanner interface{}) (*entity.Material, error) {
	var (
		id                   int
		article              string
		materialTypeID       int
		name                 string
		description          sql.NullString
		measurementUnitID    int
		packageQuantity      float64
		costPerUnit          float64
		stockQuantity        float64
		minStockQuantity     float64
		imagePath            sql.NullString
		createdAt            time.Time
		updatedAt            time.Time
	)

	var err error
	switch s := scanner.(type) {
	case *sql.Row:
		err = s.Scan(
			&id, &article, &materialTypeID, &name, &description, &measurementUnitID,
			&packageQuantity, &costPerUnit, &stockQuantity, &minStockQuantity,
			&imagePath, &createdAt, &updatedAt,
		)
	case *sql.Rows:
		err = s.Scan(
			&id, &article, &materialTypeID, &name, &description, &measurementUnitID,
			&packageQuantity, &costPerUnit, &stockQuantity, &minStockQuantity,
			&imagePath, &createdAt, &updatedAt,
		)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип scanner")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("материал не найден")
		}
		return nil, fmt.Errorf("ошибка сканирования материала: %w", err)
	}

	// Создаем материал через конструктор entity
	material, err := entity.NewMaterial(
		article,
		entity.ID(materialTypeID),
		name,
		entity.ID(measurementUnitID),
		packageQuantity,
		entity.Money(costPerUnit),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания entity материала: %w", err)
	}

	// Устанавливаем дополнительные поля
	if description.Valid {
		// material.SetDescription(&description.String)
	}

	// Примечание: В реальном проекте нужно добавить методы для установки всех полей

	return material, nil
} 