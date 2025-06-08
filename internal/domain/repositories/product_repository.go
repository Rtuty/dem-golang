package repositories

import "wallpaper-system/internal/domain/entities"

// ProductRepository определяет интерфейс для работы с продукцией
type ProductRepository interface {
	// GetAll возвращает список всей продукции
	GetAll() ([]entities.Product, error)

	// GetByID возвращает продукцию по ID
	GetByID(id int) (*entities.Product, error)

	// Create создает новую продукцию
	Create(product *entities.Product) error

	// Update обновляет продукцию
	Update(product *entities.Product) error

	// Delete удаляет продукцию
	Delete(id int) error

	// GetProductTypes возвращает все типы продукции
	GetProductTypes() ([]entities.ProductType, error)

	// GetProductTypeByID возвращает тип продукции по ID
	GetProductTypeByID(id int) (*entities.ProductType, error)

	// GetMaterialsForProduct возвращает материалы для продукции
	GetMaterialsForProduct(productID int) ([]entities.ProductMaterial, error)
}
