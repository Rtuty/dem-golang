package services

import (
	"math"

	"wallpaper-system/internal/repository"
)

// CalculatorService представляет сервис для расчетов
type CalculatorService struct {
	materialRepo *repository.MaterialRepository
	productRepo  *repository.ProductRepository
}

// NewCalculatorService создает новый сервис калькулятора
func NewCalculatorService(materialRepo *repository.MaterialRepository, productRepo *repository.ProductRepository) *CalculatorService {
	return &CalculatorService{
		materialRepo: materialRepo,
		productRepo:  productRepo,
	}
}

// CalculateRequiredMaterialAmount рассчитывает необходимое количество материала для производства
// Метод принимает следующие параметры:
// - productTypeID: идентификатор типа продукции (целое число)
// - materialTypeID: идентификатор типа материала (целое число)
// - productQuantity: количество продукции (целое число, > 0)
// - productParam1: первый параметр продукции (вещественное число, > 0)
// - productParam2: второй параметр продукции (вещественное число, > 0)
// - materialInStock: количество материала на складе (вещественное число, >= 0)
// Возвращает целое число - количество необходимого материала или -1 в случае ошибки
func (s *CalculatorService) CalculateRequiredMaterialAmount(
	productTypeID int,
	materialTypeID int,
	productQuantity int,
	productParam1 float64,
	productParam2 float64,
	materialInStock float64,
) int {
	// Проверяем входные данные
	if productTypeID <= 0 || materialTypeID <= 0 || productQuantity <= 0 ||
		productParam1 <= 0 || productParam2 <= 0 || materialInStock < 0 {
		return -1
	}

	// Получаем тип продукции
	productType, err := s.productRepo.GetProductTypeByID(productTypeID)
	if err != nil {
		return -1
	}

	// Получаем тип материала
	materialType, err := s.materialRepo.GetMaterialTypeByID(materialTypeID)
	if err != nil {
		return -1
	}

	// Рассчитываем необходимое количество материала на одну единицу продукции
	materialPerUnit := productParam1 * productParam2 * productType.Coefficient

	// Общее количество материала для всех единиц продукции
	totalMaterialNeeded := materialPerUnit * float64(productQuantity)

	// Учитываем процент брака материала
	wasteMultiplier := 1.0 + (materialType.WastePercentage / 100.0)
	totalMaterialWithWaste := totalMaterialNeeded * wasteMultiplier

	// Учитываем материал на складе
	materialToPurchase := totalMaterialWithWaste - materialInStock

	// Если материала на складе достаточно, возвращаем 0
	if materialToPurchase <= 0 {
		return 0
	}

	// Округляем до целого числа в большую сторону
	requiredQuantity := int(math.Ceil(materialToPurchase))

	return requiredQuantity
}
