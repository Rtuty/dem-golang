package repository

import (
	"context"

	"wallpaper-system/internal/domain/entity"
)

// MaterialRepository определяет интерфейс для работы с материалами
type MaterialRepository interface {
	// Create создает новый материал
	Create(ctx context.Context, material *entity.Material) error
	
	// GetByID возвращает материал по ID
	GetByID(ctx context.Context, id entity.ID) (*entity.Material, error)
	
	// GetAll возвращает все материалы с пагинацией
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Material, error)
	
	// Update обновляет материал
	Update(ctx context.Context, material *entity.Material) error
	
	// Delete удаляет материал
	Delete(ctx context.Context, id entity.ID) error
	
	// GetByArticle возвращает материал по артикулу
	GetByArticle(ctx context.Context, article string) (*entity.Material, error)
	
	// ExistsByArticle проверяет существование материала по артикулу
	ExistsByArticle(ctx context.Context, article string) (bool, error)
	
	// GetForProduct возвращает материалы для продукции
	GetForProduct(ctx context.Context, productID entity.ID) ([]*entity.Material, error)
	
	// GetLowStockMaterials возвращает материалы с низким остатком
	GetLowStockMaterials(ctx context.Context) ([]*entity.Material, error)
	
	// UpdateStock обновляет остаток материала
	UpdateStock(ctx context.Context, materialID entity.ID, newQuantity float64) error
	
	// Count возвращает общее количество материалов
	Count(ctx context.Context) (int, error)
	
	// Search ищет материалы по наименованию или артикулу
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Material, error)
}

// MaterialTypeRepository определяет интерфейс для работы с типами материалов
type MaterialTypeRepository interface {
	// Create создает новый тип материала
	Create(ctx context.Context, materialType *entity.MaterialType) error
	
	// GetByID возвращает тип материала по ID
	GetByID(ctx context.Context, id entity.ID) (*entity.MaterialType, error)
	
	// GetAll возвращает все типы материалов
	GetAll(ctx context.Context) ([]*entity.MaterialType, error)
	
	// Update обновляет тип материала
	Update(ctx context.Context, materialType *entity.MaterialType) error
	
	// Delete удаляет тип материала
	Delete(ctx context.Context, id entity.ID) error
	
	// ExistsByName проверяет существование типа материала по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// MeasurementUnitRepository определяет интерфейс для работы с единицами измерения
type MeasurementUnitRepository interface {
	// Create создает новую единицу измерения
	Create(ctx context.Context, unit *entity.MeasurementUnit) error
	
	// GetByID возвращает единицу измерения по ID
	GetByID(ctx context.Context, id entity.ID) (*entity.MeasurementUnit, error)
	
	// GetAll возвращает все единицы измерения
	GetAll(ctx context.Context) ([]*entity.MeasurementUnit, error)
	
	// Update обновляет единицу измерения
	Update(ctx context.Context, unit *entity.MeasurementUnit) error
	
	// Delete удаляет единицу измерения
	Delete(ctx context.Context, id entity.ID) error
	
	// ExistsByName проверяет существование единицы измерения по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
	
	// ExistsByAbbreviation проверяет существование единицы измерения по сокращению
	ExistsByAbbreviation(ctx context.Context, abbreviation string) (bool, error)
} 