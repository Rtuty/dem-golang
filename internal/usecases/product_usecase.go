package usecases

import (
	"fmt"
	"math"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/repositories"
)

// ProductUseCase содержит бизнес-логику для работы с продукцией
type ProductUseCase struct {
	productRepo  repositories.ProductRepository
	materialRepo repositories.MaterialRepository
}

// NewProductUseCase создает новый use case продукции
func NewProductUseCase(
	productRepo repositories.ProductRepository,
	materialRepo repositories.MaterialRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		materialRepo: materialRepo,
	}
}

// GetAllProducts возвращает список всей продукции с рассчитанными ценами
func (uc *ProductUseCase) GetAllProducts() ([]entities.Product, error) {
	products, err := uc.productRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Рассчитываем цены для каждой продукции
	for i := range products {
		price, err := uc.calculateProductPrice(&products[i])
		if err == nil && price > 0 {
			products[i].CalculatedPrice = &price
		}
	}

	return products, nil
}

// GetProductByID возвращает продукцию по ID с материалами и рассчитанной ценой
func (uc *ProductUseCase) GetProductByID(id int) (*entities.Product, error) {
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Рассчитываем цену
	price, err := uc.calculateProductPrice(product)
	if err == nil && price > 0 {
		product.CalculatedPrice = &price
	}

	return product, nil
}

// CreateProduct создает новую продукцию
func (uc *ProductUseCase) CreateProduct(product *entities.Product) error {
	// Валидация
	if err := product.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации: %w", err)
	}

	// Проверяем существование типа продукции
	_, err := uc.productRepo.GetProductTypeByID(product.ProductTypeID)
	if err != nil {
		return fmt.Errorf("тип продукции не найден: %w", err)
	}

	return uc.productRepo.Create(product)
}

// UpdateProduct обновляет продукцию
func (uc *ProductUseCase) UpdateProduct(product *entities.Product) error {
	// Валидация
	if err := product.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации: %w", err)
	}

	// Проверяем существование
	existing, err := uc.productRepo.GetByID(product.ID)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	// Проверяем тип продукции если он изменился
	if existing.ProductTypeID != product.ProductTypeID {
		_, err := uc.productRepo.GetProductTypeByID(product.ProductTypeID)
		if err != nil {
			return fmt.Errorf("тип продукции не найден: %w", err)
		}
	}

	return uc.productRepo.Update(product)
}

// DeleteProduct удаляет продукцию
func (uc *ProductUseCase) DeleteProduct(id int) error {
	// Проверяем существование
	_, err := uc.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	return uc.productRepo.Delete(id)
}

// GetProductTypes возвращает все типы продукции
func (uc *ProductUseCase) GetProductTypes() ([]entities.ProductType, error) {
	return uc.productRepo.GetProductTypes()
}

// calculateProductPrice рассчитывает стоимость продукции
func (uc *ProductUseCase) calculateProductPrice(product *entities.Product) (float64, error) {
	// Получаем тип продукции
	if product.ProductType == nil {
		productType, err := uc.productRepo.GetProductTypeByID(product.ProductTypeID)
		if err != nil {
			return 0, fmt.Errorf("ошибка получения типа продукции: %w", err)
		}
		product.ProductType = productType
	}

	// Получаем материалы если их нет
	if len(product.Materials) == 0 {
		materials, err := uc.productRepo.GetMaterialsForProduct(product.ID)
		if err != nil {
			return 0, fmt.Errorf("ошибка получения материалов: %w", err)
		}
		product.Materials = materials
	}

	// Используем доменную логику для расчета
	price := product.CalculatePrice()

	// Округляем до 2 знаков после запятой
	return math.Round(price*100) / 100, nil
}
