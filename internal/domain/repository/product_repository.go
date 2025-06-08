package repository

import (
	"context"

	"wallpaper-system/internal/domain/entity"
)

// ProductRepository определяет интерфейс для работы с продукцией
type ProductRepository interface {
	// Create создает новую продукцию
	Create(ctx context.Context, product *entity.Product) error
	
	// GetByID возвращает продукцию по ID
	GetByID(ctx context.Context, id entity.ID) (*entity.Product, error)
	
	// GetAll возвращает все продукции с пагинацией
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	
	// Update обновляет продукцию
	Update(ctx context.Context, product *entity.Product) error
	
	// Delete удаляет продукцию
	Delete(ctx context.Context, id entity.ID) error
	
	// GetByArticle возвращает продукцию по артикулу
	GetByArticle(ctx context.Context, article string) (*entity.Product, error)
	
	// ExistsByArticle проверяет существование продукции по артикулу
	ExistsByArticle(ctx context.Context, article string) (bool, error)
	
	// GetMaterials возвращает материалы для продукции
	GetMaterials(ctx context.Context, productID entity.ID) ([]entity.ProductMaterial, error)
	
	// AddMaterial добавляет материал к продукции
	AddMaterial(ctx context.Context, productMaterial *entity.ProductMaterial) error
	
	// RemoveMaterial удаляет материал из продукции
	RemoveMaterial(ctx context.Context, productID, materialID entity.ID) error
	
	// UpdateMaterialQuantity обновляет количество материала для продукции
	UpdateMaterialQuantity(ctx context.Context, productID, materialID entity.ID, quantity float64) error
	
	// Count возвращает общее количество продукции
	Count(ctx context.Context) (int, error)
	
	// GetWithCalculatedPrices возвращает продукции с рассчитанными ценами
	GetWithCalculatedPrices(ctx context.Context, limit, offset int) ([]*entity.Product, error)
}

// ProductTypeRepository определяет интерфейс для работы с типами продукции
type ProductTypeRepository interface {
	// Create создает новый тип продукции
	Create(ctx context.Context, productType *entity.ProductType) error
	
	// GetByID возвращает тип продукции по ID
	GetByID(ctx context.Context, id entity.ID) (*entity.ProductType, error)
	
	// GetAll возвращает все типы продукции
	GetAll(ctx context.Context) ([]*entity.ProductType, error)
	
	// Update обновляет тип продукции
	Update(ctx context.Context, productType *entity.ProductType) error
	
	// Delete удаляет тип продукции
	Delete(ctx context.Context, id entity.ID) error
	
	// ExistsByName проверяет существование типа продукции по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
} 