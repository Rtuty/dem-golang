package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"
)

// MaterialCalculationRequest представляет запрос на расчет материалов
type MaterialCalculationRequest struct {
	ProductTypeID   entity.ID `json:"product_type_id" validate:"required"`
	MaterialTypeID  entity.ID `json:"material_type_id" validate:"required"`
	ProductQuantity int       `json:"product_quantity" validate:"required,min=1"`
	ProductParam1   float64   `json:"product_param1" validate:"required,min=0"`
	ProductParam2   float64   `json:"product_param2" validate:"required,min=0"`
	MaterialInStock float64   `json:"material_in_stock" validate:"min=0"`
}

// MaterialCalculationResult представляет результат расчета материалов
type MaterialCalculationResult struct {
	RequiredQuantity     int     `json:"required_quantity"`
	MaterialPerUnit      float64 `json:"material_per_unit"`
	TotalMaterialNeeded  float64 `json:"total_material_needed"`
	MaterialWithWaste    float64 `json:"material_with_waste"`
	WastePercentage      float64 `json:"waste_percentage"`
	ProductCoefficient   float64 `json:"product_coefficient"`
	MaterialToPurchase   float64 `json:"material_to_purchase"`
	IsStockSufficient    bool    `json:"is_stock_sufficient"`
}

// CalculatorService представляет доменный сервис для расчета материалов
type CalculatorService struct {
	productTypeRepo  repository.ProductTypeRepository
	materialTypeRepo repository.MaterialTypeRepository
}

// NewCalculatorService создает новый сервис калькулятора
func NewCalculatorService(
	productTypeRepo repository.ProductTypeRepository,
	materialTypeRepo repository.MaterialTypeRepository,
) *CalculatorService {
	return &CalculatorService{
		productTypeRepo:  productTypeRepo,
		materialTypeRepo: materialTypeRepo,
	}
}

// CalculateRequiredMaterial рассчитывает необходимое количество материала для производства
func (s *CalculatorService) CalculateRequiredMaterial(
	ctx context.Context,
	req *MaterialCalculationRequest,
) (*MaterialCalculationResult, error) {
	// Валидация входных данных
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("ошибка валидации запроса: %w", err)
	}

	// Получаем тип продукции
	productType, err := s.productTypeRepo.GetByID(ctx, req.ProductTypeID)
	if err != nil {
		return nil, fmt.Errorf("тип продукции не найден: %w", err)
	}

	// Получаем тип материала
	materialType, err := s.materialTypeRepo.GetByID(ctx, req.MaterialTypeID)
	if err != nil {
		return nil, fmt.Errorf("тип материала не найден: %w", err)
	}

	// Выполняем расчеты
	result := s.performCalculation(req, productType, materialType)

	return result, nil
}

// validateRequest проверяет корректность запроса
func (s *CalculatorService) validateRequest(req *MaterialCalculationRequest) error {
	if req.ProductQuantity <= 0 {
		return errors.New("количество продукции должно быть положительным")
	}
	if req.ProductParam1 <= 0 {
		return errors.New("первый параметр продукции должен быть положительным")
	}
	if req.ProductParam2 <= 0 {
		return errors.New("второй параметр продукции должен быть положительным")
	}
	if req.MaterialInStock < 0 {
		return errors.New("количество материала на складе не может быть отрицательным")
	}
	return nil
}

// performCalculation выполняет основные расчеты
func (s *CalculatorService) performCalculation(
	req *MaterialCalculationRequest,
	productType *entity.ProductType,
	materialType *entity.MaterialType,
) *MaterialCalculationResult {
	// Рассчитываем необходимое количество материала на одну единицу продукции
	materialPerUnit := req.ProductParam1 * req.ProductParam2 * productType.Coefficient()

	// Общее количество материала для всех единиц продукции
	totalMaterialNeeded := materialPerUnit * float64(req.ProductQuantity)

	// Учитываем процент брака материала
	wasteMultiplier := 1.0 + (materialType.WastePercentage() / 100.0)
	materialWithWaste := totalMaterialNeeded * wasteMultiplier

	// Учитываем материал на складе
	materialToPurchase := materialWithWaste - req.MaterialInStock

	// Определяем, достаточно ли материала на складе
	isStockSufficient := materialToPurchase <= 0

	// Округляем до целого числа в большую сторону для покупки
	requiredQuantity := 0
	if !isStockSufficient {
		requiredQuantity = int(math.Ceil(materialToPurchase))
	}

	return &MaterialCalculationResult{
		RequiredQuantity:     requiredQuantity,
		MaterialPerUnit:      materialPerUnit,
		TotalMaterialNeeded:  totalMaterialNeeded,
		MaterialWithWaste:    materialWithWaste,
		WastePercentage:      materialType.WastePercentage(),
		ProductCoefficient:   productType.Coefficient(),
		MaterialToPurchase:   math.Max(0, materialToPurchase),
		IsStockSufficient:    isStockSufficient,
	}
}

// CalculateProductionCost рассчитывает стоимость производства
func (s *CalculatorService) CalculateProductionCost(
	ctx context.Context,
	productID entity.ID,
	quantity int,
	materialRepo repository.MaterialRepository,
	productRepo repository.ProductRepository,
) (*ProductionCostResult, error) {
	if quantity <= 0 {
		return nil, errors.New("количество должно быть положительным")
	}

	// Получаем продукцию
	product, err := productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("продукция не найдена: %w", err)
	}

	// Получаем материалы для продукции
	materials, err := productRepo.GetMaterials(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов: %w", err)
	}

	if len(materials) == 0 {
		return nil, errors.New("у продукции нет связанных материалов")
	}

	var totalCost entity.Money
	var materialCosts []MaterialCostDetail

	// Рассчитываем стоимость каждого материала
	for _, pm := range materials {
		if pm.Material() == nil {
			continue
		}

		material := pm.Material()
		costPerUnit := material.CostPerUnit()
		quantityNeeded := pm.QuantityPerUnit() * float64(quantity)
		totalMaterialCost := entity.Money(quantityNeeded) * costPerUnit

		materialCosts = append(materialCosts, MaterialCostDetail{
			MaterialID:    material.ID(),
			MaterialName:  material.Name(),
			CostPerUnit:   costPerUnit,
			QuantityNeeded: quantityNeeded,
			TotalCost:     totalMaterialCost,
		})

		totalCost += totalMaterialCost
	}

	return &ProductionCostResult{
		ProductID:      productID,
		ProductName:    product.Name(),
		Quantity:       quantity,
		TotalCost:      totalCost,
		CostPerUnit:    totalCost / entity.Money(quantity),
		MaterialCosts:  materialCosts,
	}, nil
}

// ProductionCostResult представляет результат расчета стоимости производства
type ProductionCostResult struct {
	ProductID      entity.ID            `json:"product_id"`
	ProductName    string               `json:"product_name"`
	Quantity       int                  `json:"quantity"`
	TotalCost      entity.Money         `json:"total_cost"`
	CostPerUnit    entity.Money         `json:"cost_per_unit"`
	MaterialCosts  []MaterialCostDetail `json:"material_costs"`
}

// MaterialCostDetail представляет детали стоимости материала
type MaterialCostDetail struct {
	MaterialID     entity.ID    `json:"material_id"`
	MaterialName   string       `json:"material_name"`
	CostPerUnit    entity.Money `json:"cost_per_unit"`
	QuantityNeeded float64      `json:"quantity_needed"`
	TotalCost      entity.Money `json:"total_cost"`
}

// OptimizeMaterialUsage оптимизирует использование материалов
func (s *CalculatorService) OptimizeMaterialUsage(
	ctx context.Context,
	productID entity.ID,
	targetQuantity int,
	materialRepo repository.MaterialRepository,
	productRepo repository.ProductRepository,
) (*MaterialOptimizationResult, error) {
	// Получаем материалы для продукции
	materials, err := productRepo.GetMaterials(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов: %w", err)
	}

	var recommendations []MaterialRecommendation

	for _, pm := range materials {
		if pm.Material() == nil {
			continue
		}

		material := pm.Material()
		requiredQuantity := pm.QuantityPerUnit() * float64(targetQuantity)
		
		recommendation := MaterialRecommendation{
			MaterialID:       material.ID(),
			MaterialName:     material.Name(),
			RequiredQuantity: requiredQuantity,
			CurrentStock:     material.StockQuantity(),
			MinStock:         material.MinStockQuantity(),
			IsLowStock:       material.IsLowStock(),
		}

		if material.StockQuantity() < requiredQuantity {
			recommendation.NeedToPurchase = requiredQuantity - material.StockQuantity()
			recommendation.RecommendedPurchase = math.Ceil(recommendation.NeedToPurchase / material.PackageQuantity()) * material.PackageQuantity()
		}

		recommendations = append(recommendations, recommendation)
	}

	return &MaterialOptimizationResult{
		ProductID:       productID,
		TargetQuantity:  targetQuantity,
		Recommendations: recommendations,
	}, nil
}

// MaterialOptimizationResult представляет результат оптимизации материалов
type MaterialOptimizationResult struct {
	ProductID       entity.ID                `json:"product_id"`
	TargetQuantity  int                      `json:"target_quantity"`
	Recommendations []MaterialRecommendation `json:"recommendations"`
}

// MaterialRecommendation представляет рекомендацию по материалу
type MaterialRecommendation struct {
	MaterialID           entity.ID `json:"material_id"`
	MaterialName         string    `json:"material_name"`
	RequiredQuantity     float64   `json:"required_quantity"`
	CurrentStock         float64   `json:"current_stock"`
	MinStock             float64   `json:"min_stock"`
	NeedToPurchase       float64   `json:"need_to_purchase"`
	RecommendedPurchase  float64   `json:"recommended_purchase"`
	IsLowStock           bool      `json:"is_low_stock"`
} 