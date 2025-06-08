package services

import (
	"fmt"
	"math"

	"wallpaper-system/internal/models"
	"wallpaper-system/internal/repository"
)

// MaterialService представляет сервис для работы с материалами
type MaterialService struct {
	materialRepo *repository.MaterialRepository
}

// NewMaterialService создает новый сервис материалов
func NewMaterialService(materialRepo *repository.MaterialRepository) *MaterialService {
	return &MaterialService{
		materialRepo: materialRepo,
	}
}

// GetAllMaterials возвращает все материалы
func (s *MaterialService) GetAllMaterials() ([]models.Material, error) {
	materials, err := s.materialRepo.GetAllMaterials()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка материалов: %w", err)
	}

	return materials, nil
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (s *MaterialService) GetMaterialsForProduct(productID int) ([]models.MaterialForProduct, error) {
	materials, err := s.materialRepo.GetMaterialsForProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов для продукции: %w", err)
	}

	return materials, nil
}

// CalculateRequiredMaterial рассчитывает необходимое количество материала для производства
func (s *MaterialService) CalculateRequiredMaterial(req *models.MaterialCalculationRequest) (*models.MaterialCalculationResponse, error) {
	// Проверяем существование типа продукции
	productType, err := s.materialRepo.GetProductTypeByID(req.ProductTypeID)
	if err != nil {
		return &models.MaterialCalculationResponse{RequiredQuantity: -1}, nil
	}

	// Проверяем существование типа материала
	materialType, err := s.materialRepo.GetMaterialTypeByID(req.MaterialTypeID)
	if err != nil {
		return &models.MaterialCalculationResponse{RequiredQuantity: -1}, nil
	}

	// Валидация входных данных
	if req.ProductQuantity <= 0 || req.ProductParam1 <= 0 || req.ProductParam2 <= 0 || req.MaterialInStock < 0 {
		return &models.MaterialCalculationResponse{RequiredQuantity: -1}, nil
	}

	// Рассчитываем необходимое количество материала на одну единицу продукции
	materialPerUnit := req.ProductParam1 * req.ProductParam2 * productType.Coefficient

	// Общее количество материала для всех единиц продукции
	totalMaterialNeeded := materialPerUnit * float64(req.ProductQuantity)

	// Учитываем процент брака материала
	wasteMultiplier := 1.0 + (materialType.WastePercentage / 100.0)
	totalMaterialWithWaste := totalMaterialNeeded * wasteMultiplier

	// Учитываем материал на складе
	materialToPurchase := totalMaterialWithWaste - req.MaterialInStock

	// Если материала на складе достаточно, возвращаем 0
	if materialToPurchase <= 0 {
		return &models.MaterialCalculationResponse{RequiredQuantity: 0}, nil
	}

	// Округляем до целого числа в большую сторону
	requiredQuantity := int(math.Ceil(materialToPurchase))

	return &models.MaterialCalculationResponse{RequiredQuantity: requiredQuantity}, nil
}

// GetMaterialTypes возвращает все типы материалов
func (s *MaterialService) GetMaterialTypes() ([]models.MaterialType, error) {
	types, err := s.materialRepo.GetMaterialTypes()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения типов материалов: %w", err)
	}

	return types, nil
}
