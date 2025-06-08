package usecases

import "wallpaper-system/internal/domain/entities"

// ProductUseCaseInterface определяет интерфейс для работы с продукцией
type ProductUseCaseInterface interface {
	GetAllProducts() ([]entities.Product, error)
	GetProductByID(id int) (*entities.Product, error)
	CreateProduct(product *entities.Product) error
	UpdateProduct(product *entities.Product) error
	DeleteProduct(id int) error
	GetProductTypes() ([]entities.ProductType, error)
}

// MaterialUseCaseInterface определяет интерфейс для работы с материалами
type MaterialUseCaseInterface interface {
	GetAllMaterials() ([]entities.Material, error)
	GetMaterialByID(id int) (*entities.Material, error)
	CreateMaterial(material *entities.Material) error
	UpdateMaterial(material *entities.Material) error
	DeleteMaterial(id int) error
	GetMaterialTypes() ([]entities.MaterialType, error)
	GetMeasurementUnits() ([]entities.MeasurementUnit, error)
	GetMaterialsForProduct(productID int) ([]entities.Material, error)
}

// CalculatorUseCaseInterface определяет интерфейс для калькулятора
type CalculatorUseCaseInterface interface {
	CalculateRequiredMaterial(request *entities.MaterialCalculationRequest) (int, error)
}
