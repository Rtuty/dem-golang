package usecase

import (
	"context"
	"fmt"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/service"
)

// CalculatorUseCase представляет вариант использования для калькулятора материалов
type CalculatorUseCase struct {
	calculatorService *service.CalculatorService
}

// NewCalculatorUseCase создает новый use case для калькулятора
func NewCalculatorUseCase(calculatorService *service.CalculatorService) *CalculatorUseCase {
	return &CalculatorUseCase{
		calculatorService: calculatorService,
	}
}

// MaterialCalculationRequest представляет запрос на расчет материалов
type MaterialCalculationRequest struct {
	ProductTypeID   int     `json:"product_type_id" validate:"required,min=1"`
	MaterialTypeID  int     `json:"material_type_id" validate:"required,min=1"`
	ProductQuantity int     `json:"product_quantity" validate:"required,min=1"`
	ProductParam1   float64 `json:"product_param1" validate:"required,min=0"`
	ProductParam2   float64 `json:"product_param2" validate:"required,min=0"`
	MaterialInStock float64 `json:"material_in_stock" validate:"min=0"`
}

// MaterialCalculationResponse представляет ответ с результатами расчета
type MaterialCalculationResponse struct {
	RequiredQuantity     int     `json:"required_quantity"`
	MaterialPerUnit      float64 `json:"material_per_unit"`
	TotalMaterialNeeded  float64 `json:"total_material_needed"`
	MaterialWithWaste    float64 `json:"material_with_waste"`
	WastePercentage      float64 `json:"waste_percentage"`
	ProductCoefficient   float64 `json:"product_coefficient"`
	MaterialToPurchase   float64 `json:"material_to_purchase"`
	IsStockSufficient    bool    `json:"is_stock_sufficient"`
	Message              string  `json:"message"`
}

// ProductionCostRequest представляет запрос на расчет стоимости производства
type ProductionCostRequest struct {
	ProductID int `json:"product_id" validate:"required,min=1"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

// ProductionCostResponse представляет ответ с расчетом стоимости производства
type ProductionCostResponse struct {
	ProductID      int                         `json:"product_id"`
	ProductName    string                      `json:"product_name"`
	Quantity       int                         `json:"quantity"`
	TotalCost      float64                     `json:"total_cost"`
	CostPerUnit    float64                     `json:"cost_per_unit"`
	MaterialCosts  []MaterialCostDetailResponse `json:"material_costs"`
}

// MaterialCostDetailResponse представляет детали стоимости материала
type MaterialCostDetailResponse struct {
	MaterialID     int     `json:"material_id"`
	MaterialName   string  `json:"material_name"`
	CostPerUnit    float64 `json:"cost_per_unit"`
	QuantityNeeded float64 `json:"quantity_needed"`
	TotalCost      float64 `json:"total_cost"`
}

// MaterialOptimizationRequest представляет запрос на оптимизацию материалов
type MaterialOptimizationRequest struct {
	ProductID      int `json:"product_id" validate:"required,min=1"`
	TargetQuantity int `json:"target_quantity" validate:"required,min=1"`
}

// MaterialOptimizationResponse представляет ответ с рекомендациями по материалам
type MaterialOptimizationResponse struct {
	ProductID       int                               `json:"product_id"`
	TargetQuantity  int                               `json:"target_quantity"`
	Recommendations []MaterialRecommendationResponse  `json:"recommendations"`
}

// MaterialRecommendationResponse представляет рекомендацию по материалу
type MaterialRecommendationResponse struct {
	MaterialID           int     `json:"material_id"`
	MaterialName         string  `json:"material_name"`
	RequiredQuantity     float64 `json:"required_quantity"`
	CurrentStock         float64 `json:"current_stock"`
	MinStock             float64 `json:"min_stock"`
	NeedToPurchase       float64 `json:"need_to_purchase"`
	RecommendedPurchase  float64 `json:"recommended_purchase"`
	IsLowStock           bool    `json:"is_low_stock"`
	Status               string  `json:"status"`
}

// CalculateRequiredMaterial рассчитывает необходимое количество материала
func (uc *CalculatorUseCase) CalculateRequiredMaterial(
	ctx context.Context,
	req *MaterialCalculationRequest,
) (*MaterialCalculationResponse, error) {
	// Преобразуем запрос в доменный объект
	domainReq := &service.MaterialCalculationRequest{
		ProductTypeID:   entity.ID(req.ProductTypeID),
		MaterialTypeID:  entity.ID(req.MaterialTypeID),
		ProductQuantity: req.ProductQuantity,
		ProductParam1:   req.ProductParam1,
		ProductParam2:   req.ProductParam2,
		MaterialInStock: req.MaterialInStock,
	}

	// Выполняем расчет
	result, err := uc.calculatorService.CalculateRequiredMaterial(ctx, domainReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка расчета материалов: %w", err)
	}

	// Формируем сообщение для пользователя
	message := uc.generateMessage(result)

	// Преобразуем результат в ответ
	response := &MaterialCalculationResponse{
		RequiredQuantity:     result.RequiredQuantity,
		MaterialPerUnit:      result.MaterialPerUnit,
		TotalMaterialNeeded:  result.TotalMaterialNeeded,
		MaterialWithWaste:    result.MaterialWithWaste,
		WastePercentage:      result.WastePercentage,
		ProductCoefficient:   result.ProductCoefficient,
		MaterialToPurchase:   result.MaterialToPurchase,
		IsStockSufficient:    result.IsStockSufficient,
		Message:              message,
	}

	return response, nil
}

// generateMessage генерирует понятное для пользователя сообщение
func (uc *CalculatorUseCase) generateMessage(result *service.MaterialCalculationResult) string {
	if result.IsStockSufficient {
		return "Материала на складе достаточно для производства"
	}

	if result.RequiredQuantity > 0 {
		return fmt.Sprintf("Необходимо закупить %d единиц материала", result.RequiredQuantity)
	}

	return "Расчет выполнен"
}

// CalculateProductionCost рассчитывает стоимость производства продукции
func (uc *CalculatorUseCase) CalculateProductionCost(
	ctx context.Context,
	req *ProductionCostRequest,
	materialRepo interface{}, // Будет заменено на конкретную реализацию
	productRepo interface{},  // Будет заменено на конкретную реализацию
) (*ProductionCostResponse, error) {
	// Пока что возвращаем заглушку, так как нужны конкретные репозитории
	// В полной реализации здесь будет вызов calculatorService.CalculateProductionCost
	return &ProductionCostResponse{
		ProductID:   req.ProductID,
		ProductName: "Заглушка",
		Quantity:    req.Quantity,
		TotalCost:   0,
		CostPerUnit: 0,
		MaterialCosts: []MaterialCostDetailResponse{},
	}, nil
}

// OptimizeMaterialUsage оптимизирует использование материалов
func (uc *CalculatorUseCase) OptimizeMaterialUsage(
	ctx context.Context,
	req *MaterialOptimizationRequest,
	materialRepo interface{}, // Будет заменено на конкретную реализацию
	productRepo interface{},  // Будет заменено на конкретную реализацию
) (*MaterialOptimizationResponse, error) {
	// Пока что возвращаем заглушку
	// В полной реализации здесь будет вызов calculatorService.OptimizeMaterialUsage
	return &MaterialOptimizationResponse{
		ProductID:      req.ProductID,
		TargetQuantity: req.TargetQuantity,
		Recommendations: []MaterialRecommendationResponse{},
	}, nil
}

// ValidateCalculationRequest проверяет корректность запроса на расчет
func (uc *CalculatorUseCase) ValidateCalculationRequest(req *MaterialCalculationRequest) []string {
	var errors []string

	if req.ProductTypeID <= 0 {
		errors = append(errors, "ID типа продукции должен быть положительным")
	}

	if req.MaterialTypeID <= 0 {
		errors = append(errors, "ID типа материала должен быть положительным")
	}

	if req.ProductQuantity <= 0 {
		errors = append(errors, "Количество продукции должно быть положительным")
	}

	if req.ProductParam1 <= 0 {
		errors = append(errors, "Первый параметр продукции должен быть положительным")
	}

	if req.ProductParam2 <= 0 {
		errors = append(errors, "Второй параметр продукции должен быть положительным")
	}

	if req.MaterialInStock < 0 {
		errors = append(errors, "Количество материала на складе не может быть отрицательным")
	}

	return errors
}

// GetCalculationSummary возвращает сводку по расчету
func (uc *CalculatorUseCase) GetCalculationSummary(
	ctx context.Context,
	req *MaterialCalculationRequest,
) (*CalculationSummaryResponse, error) {
	// Выполняем основной расчет
	result, err := uc.CalculateRequiredMaterial(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сводки расчета: %w", err)
	}

	// Формируем сводную информацию
	summary := &CalculationSummaryResponse{
		Calculation:    *result,
		Efficiency:     uc.calculateEfficiency(result),
		Recommendations: uc.generateRecommendations(result),
	}

	return summary, nil
}

// CalculationSummaryResponse представляет сводку по расчету
type CalculationSummaryResponse struct {
	Calculation     MaterialCalculationResponse `json:"calculation"`
	Efficiency      float64                     `json:"efficiency"`
	Recommendations []string                    `json:"recommendations"`
}

// calculateEfficiency рассчитывает эффективность использования материала
func (uc *CalculatorUseCase) calculateEfficiency(result *MaterialCalculationResponse) float64 {
	if result.MaterialWithWaste == 0 {
		return 100.0
	}

	efficiency := (result.TotalMaterialNeeded / result.MaterialWithWaste) * 100
	return efficiency
}

// generateRecommendations генерирует рекомендации по оптимизации
func (uc *CalculatorUseCase) generateRecommendations(result *MaterialCalculationResponse) []string {
	var recommendations []string

	if result.WastePercentage > 10 {
		recommendations = append(recommendations, "Рассмотрите возможность снижения процента брака материала")
	}

	if !result.IsStockSufficient && result.RequiredQuantity > 100 {
		recommendations = append(recommendations, "Большой объем закупки - рассмотрите оптовые скидки")
	}

	if result.IsStockSufficient {
		recommendations = append(recommendations, "Материала достаточно - можно начинать производство")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Расчет выполнен корректно")
	}

	return recommendations
} 