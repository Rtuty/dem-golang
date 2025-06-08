package services

import (
	"fmt"
	"math"

	"wallpaper-system/internal/models"
	"wallpaper-system/internal/repository"
)

// ProductService представляет сервис для работы с продукцией
type ProductService struct {
	productRepo  *repository.ProductRepository
	materialRepo *repository.MaterialRepository
}

// NewProductService создает новый сервис продукции
func NewProductService(productRepo *repository.ProductRepository, materialRepo *repository.MaterialRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		materialRepo: materialRepo,
	}
}

// GetAllProducts возвращает список всей продукции с рассчитанными стоимостями
func (s *ProductService) GetAllProducts() ([]models.ProductListItem, error) {
	products, err := s.productRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка продукции: %w", err)
	}

	// Рассчитываем стоимость для каждого продукта
	for i := range products {
		calculatedPrice, err := s.CalculateProductPrice(products[i].ID)
		if err != nil {
			// Если не удалось рассчитать стоимость, продолжаем без неё
			continue
		}
		products[i].CalculatedPrice = &calculatedPrice
	}

	return products, nil
}

// GetProductByID возвращает продукцию по ID с рассчитанной стоимостью
func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Рассчитываем стоимость продукции
	calculatedPrice, err := s.CalculateProductPrice(id)
	if err != nil {
		// Если не удалось рассчитать стоимость, возвращаем продукт без расчетной стоимости
		fmt.Printf("Предупреждение: не удалось рассчитать стоимость для продукта %d: %v\n", id, err)
	} else {
		product.CalculatedPrice = &calculatedPrice
	}

	return product, nil
}

// CalculateProductPrice рассчитывает стоимость продукции на основе используемых материалов
func (s *ProductService) CalculateProductPrice(productID int) (float64, error) {
	// Получаем материалы для продукции
	materials, err := s.materialRepo.GetMaterialsForProduct(productID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения материалов для продукции: %w", err)
	}

	if len(materials) == 0 {
		return 0, fmt.Errorf("у продукции нет связанных материалов")
	}

	var totalCost float64

	// Суммируем стоимость всех материалов
	for _, material := range materials {
		materialCost := material.QuantityPerUnit * material.CostPerUnit
		totalCost += materialCost
	}

	// Получаем информацию о типе продукции для учета коэффициента
	product, err := s.productRepo.GetByID(productID)
	if err == nil && product.ProductType != nil {
		// Применяем коэффициент типа продукции
		totalCost *= product.ProductType.Coefficient
	}

	// Добавляем наценку (например, 20%)
	totalCost *= 1.2

	// Округляем до сотых
	totalCost = math.Round(totalCost*100) / 100

	// Проверяем, что стоимость не отрицательная
	if totalCost < 0 {
		totalCost = 0
	}

	return totalCost, nil
}

// CreateProduct создает новую продукцию
func (s *ProductService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	// Проверяем, что артикул уникален
	// (В реальном проекте здесь была бы проверка уникальности)

	product, err := s.productRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания продукции: %w", err)
	}

	return product, nil
}

// UpdateProduct обновляет продукцию
func (s *ProductService) UpdateProduct(id int, req *models.UpdateProductRequest) error {
	err := s.productRepo.Update(id, req)
	if err != nil {
		return fmt.Errorf("ошибка обновления продукции: %w", err)
	}

	return nil
}

// DeleteProduct удаляет продукцию
func (s *ProductService) DeleteProduct(id int) error {
	err := s.productRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("ошибка удаления продукции: %w", err)
	}

	return nil
}

// GetProductTypes возвращает все типы продукции
func (s *ProductService) GetProductTypes() ([]models.ProductType, error) {
	types, err := s.productRepo.GetProductTypes()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения типов продукции: %w", err)
	}

	return types, nil
}

// GetProductMaterials возвращает список материалов для конкретной продукции
func (s *ProductService) GetProductMaterials(productID int) ([]models.MaterialForProduct, error) {
	materials, err := s.materialRepo.GetMaterialsForProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов для продукции: %w", err)
	}

	return materials, nil
}
