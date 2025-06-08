package repositories

import "wallpaper-system/internal/domain/entities"

// MaterialRepository определяет интерфейс для работы с материалами
type MaterialRepository interface {
	// GetAll возвращает список всех материалов
	GetAll() ([]entities.Material, error)

	// GetByID возвращает материал по ID
	GetByID(id int) (*entities.Material, error)

	// GetMaterialTypeByID возвращает тип материала по ID
	GetMaterialTypeByID(id int) (*entities.MaterialType, error)

	// GetProductTypeByID возвращает тип продукции по ID
	GetProductTypeByID(id int) (*entities.ProductType, error)

	// GetMaterialTypes возвращает все типы материалов
	GetMaterialTypes() ([]entities.MaterialType, error)

	// GetMaterialsForProduct возвращает материалы для конкретной продукции
	GetMaterialsForProduct(productID int) ([]entities.Material, error)
}
