package service

import (
	"fmt"
	"math"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"
)

// CalculationService предоставляет методы для расчетов
type CalculationService struct {
	productRepo      repository.ProductRepository
	materialRepo     repository.MaterialRepository
	productTypeRepo  repository.ProductTypeRepository
	materialTypeRepo repository.MaterialTypeRepository
}

// NewCalculationService создает новый сервис расчетов
func NewCalculationService(
	productRepo repository.ProductRepository,
	materialRepo repository.MaterialRepository,
	productTypeRepo repository.ProductTypeRepository,
	materialTypeRepo repository.MaterialTypeRepository,
) *CalculationService {
	return &CalculationService{
		productRepo:      productRepo,
		materialRepo:     materialRepo,
		productTypeRepo:  productTypeRepo,
		materialTypeRepo: materialTypeRepo,
	}
}

// ProductCostCalculation результат расчета стоимости продукции
type ProductCostCalculation struct {
	ProductID       entity.ID           `json:"product_id"`
	MaterialCosts   []MaterialCostItem  `json:"material_costs"`
	TotalCost       entity.Money        `json:"total_cost"`
	CostPerUnit     entity.Money        `json:"cost_per_unit"`
}

// MaterialCostItem стоимость материала в продукции
type MaterialCostItem struct {
	MaterialID     entity.ID    `json:"material_id"`
	MaterialName   string       `json:"material_name"`
	QuantityNeeded float64      `json:"quantity_needed"`
	UnitCost       entity.Money `json:"unit_cost"`
	TotalCost      entity.Money `json:"total_cost"`
}

// MaterialRequirementCalculation результат расчета потребности в материалах
type MaterialRequirementCalculation struct {
	MaterialID          entity.ID `json:"material_id"`
	MaterialName        string    `json:"material_name"`
	BaseQuantityNeeded  float64   `json:"base_quantity_needed"`
	QuantityWithDefect  float64   `json:"quantity_with_defect"`
	StockQuantity       float64   `json:"stock_quantity"`
	QuantityToPurchase  int       `json:"quantity_to_purchase"`
}

// CalculateProductCost рассчитывает стоимость продукции исходя из материалов
func (cs *CalculationService) CalculateProductCost(productID entity.ID, quantity int) (*ProductCostCalculation, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("количество продукции должно быть положительным")
	}

	// Получаем продукт
	product, err := cs.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Получаем материалы для продукции
	productMaterials, err := cs.productRepo.GetMaterials(productID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов продукции: %w", err)
	}

	if len(productMaterials) == 0 {
		return &ProductCostCalculation{
			ProductID:   productID,
			TotalCost:   0,
			CostPerUnit: 0,
		}, nil
	}

	var materialCosts []MaterialCostItem
	var totalCost entity.Money

	for _, pm := range productMaterials {
		material, err := cs.materialRepo.GetByID(pm.MaterialID())
		if err != nil {
			continue // Пропускаем недоступные материалы
		}

		quantityNeeded := pm.QuantityPerUnit() * float64(quantity)
		materialTotalCost := entity.Money(quantityNeeded) * material.CostPerUnit()

		materialCosts = append(materialCosts, MaterialCostItem{
			MaterialID:     material.ID(),
			MaterialName:   material.Name(),
			QuantityNeeded: quantityNeeded,
			UnitCost:       material.CostPerUnit(),
			TotalCost:      materialTotalCost,
		})

		totalCost += materialTotalCost
	}

	costPerUnit := totalCost / entity.Money(quantity)

	return &ProductCostCalculation{
		ProductID:     productID,
		MaterialCosts: materialCosts,
		TotalCost:     totalCost,
		CostPerUnit:   costPerUnit,
	}, nil
}

// CalculateRequiredMaterial рассчитывает количество материала с учетом брака
// Согласно техническому заданию модуля 4
func (cs *CalculationService) CalculateRequiredMaterial(
	productTypeID entity.ID,
	materialTypeID entity.ID,
	productQuantity int,
	param1 float64,
	param2 float64,
	stockQuantity float64,
) (int, error) {
	// Валидация входных параметров
	if productQuantity <= 0 {
		return -1, fmt.Errorf("количество продукции должно быть положительным")
	}

	if param1 <= 0 || param2 <= 0 {
		return -1, fmt.Errorf("параметры продукции должны быть положительными")
	}

	if stockQuantity < 0 {
		return -1, fmt.Errorf("количество на складе не может быть отрицательным")
	}

	// Получаем тип продукции
	productType, err := cs.productTypeRepo.GetByID(productTypeID)
	if err != nil {
		return -1, fmt.Errorf("неизвестный тип продукции")
	}

	// Получаем тип материала
	materialType, err := cs.materialTypeRepo.GetByID(materialTypeID)
	if err != nil {
		return -1, fmt.Errorf("неизвестный тип материала")
	}

	// Расчет количества материала на одну единицу продукции
	// Формула: param1 * param2 * коэффициент_типа_продукции
	quantityPerUnit := param1 * param2 * float64(productType.Coefficient())

	// Общее количество материала без учета брака
	baseQuantityNeeded := quantityPerUnit * float64(productQuantity)

	// Увеличиваем количество с учетом процента брака
	quantityWithDefect := materialType.CalculateWithDefect(baseQuantityNeeded)

	// Учитываем материал на складе
	quantityToPurchase := quantityWithDefect - stockQuantity

	// Если на складе достаточно материала, возвращаем 0
	if quantityToPurchase <= 0 {
		return 0, nil
	}

	// Возвращаем целое количество (округляем вверх)
	return int(math.Ceil(quantityToPurchase)), nil
}

// CalculateDetailedMaterialRequirement возвращает подробную информацию о расчете
func (cs *CalculationService) CalculateDetailedMaterialRequirement(
	productTypeID entity.ID,
	materialTypeID entity.ID,
	productQuantity int,
	param1 float64,
	param2 float64,
	stockQuantity float64,
) (*MaterialRequirementCalculation, error) {
	// Получаем количество для покупки
	quantityToPurchase, err := cs.CalculateRequiredMaterial(
		productTypeID, materialTypeID, productQuantity, param1, param2, stockQuantity,
	)
	if err != nil {
		return nil, err
	}

	// Получаем название материала
	materialType, err := cs.materialTypeRepo.GetByID(materialTypeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения типа материала: %w", err)
	}

	// Получаем тип продукции для расчета
	productType, err := cs.productTypeRepo.GetByID(productTypeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения типа продукции: %w", err)
	}

	// Повторяем расчеты для получения промежуточных значений
	quantityPerUnit := param1 * param2 * float64(productType.Coefficient())
	baseQuantityNeeded := quantityPerUnit * float64(productQuantity)
	quantityWithDefect := materialType.CalculateWithDefect(baseQuantityNeeded)

	return &MaterialRequirementCalculation{
		MaterialID:          materialTypeID,
		MaterialName:        materialType.Name(),
		BaseQuantityNeeded:  baseQuantityNeeded,
		QuantityWithDefect:  quantityWithDefect,
		StockQuantity:       stockQuantity,
		QuantityToPurchase:  quantityToPurchase,
	}, nil
}

// CalculateProductPrice рассчитывает финальную цену продукции
// с учетом себестоимости материалов и наценки
func (cs *CalculationService) CalculateProductPrice(productID entity.ID, markup float64) (entity.Money, error) {
	if markup < 0 {
		return 0, fmt.Errorf("наценка не может быть отрицательной")
	}

	// Рассчитываем стоимость материалов для одной единицы
	costCalc, err := cs.CalculateProductCost(productID, 1)
	if err != nil {
		return 0, fmt.Errorf("ошибка расчета стоимости материалов: %w", err)
	}

	// Применяем наценку
	finalPrice := costCalc.CostPerUnit * entity.Money(1.0+markup)

	// Округляем до копеек
	return entity.Money(math.Round(float64(finalPrice)*100) / 100), nil
} 