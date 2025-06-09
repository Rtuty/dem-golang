package usecases

import (
	"fmt"
	"strconv"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/repositories"
)

// MaterialUseCase содержит бизнес-логику для работы с материалами
type MaterialUseCase struct {
	materialRepo repositories.MaterialRepository
}

// NewMaterialUseCase создает новый use case материалов
func NewMaterialUseCase(materialRepo repositories.MaterialRepository) *MaterialUseCase {
	return &MaterialUseCase{
		materialRepo: materialRepo,
	}
}

// GetAllMaterials возвращает список всех материалов
func (uc *MaterialUseCase) GetAllMaterials() ([]entities.Material, error) {
	return uc.materialRepo.GetAll()
}

// GetMaterialByID возвращает материал по ID
func (uc *MaterialUseCase) GetMaterialByID(id int) (*entities.Material, error) {
	return uc.materialRepo.GetByID(id)
}

// GetMaterialTypes возвращает все типы материалов
func (uc *MaterialUseCase) GetMaterialTypes() ([]entities.MaterialType, error) {
	return uc.materialRepo.GetMaterialTypes()
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (uc *MaterialUseCase) GetMaterialsForProduct(productID int) ([]entities.Material, error) {
	return uc.materialRepo.GetMaterialsForProduct(productID)
}

// CreateMaterial создает новый материал
func (uc *MaterialUseCase) CreateMaterial(material *entities.Material) error {
	// Валидация
	if err := uc.validateMaterial(material); err != nil {
		return fmt.Errorf("ошибка валидации материала: %w", err)
	}

	return uc.materialRepo.Create(material)
}

// UpdateMaterial обновляет существующий материал
func (uc *MaterialUseCase) UpdateMaterial(material *entities.Material) error {
	// Проверяем, что материал существует
	existing, err := uc.materialRepo.GetByID(material.ID)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}
	if existing == nil {
		return entities.NewNotFoundError("материал", strconv.Itoa(material.ID))
	}

	// Валидация
	if err := uc.validateMaterial(material); err != nil {
		return fmt.Errorf("ошибка валидации материала: %w", err)
	}

	return uc.materialRepo.Update(material)
}

// DeleteMaterial удаляет материал
func (uc *MaterialUseCase) DeleteMaterial(id int) error {
	// Проверяем, что материал существует
	existing, err := uc.materialRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}
	if existing == nil {
		return entities.NewNotFoundError("материал", strconv.Itoa(id))
	}

	return uc.materialRepo.Delete(id)
}

// GetMeasurementUnits возвращает все единицы измерения
func (uc *MaterialUseCase) GetMeasurementUnits() ([]entities.MeasurementUnit, error) {
	return uc.materialRepo.GetMeasurementUnits()
}

// validateMaterial проверяет корректность данных материала
func (uc *MaterialUseCase) validateMaterial(material *entities.Material) error {
	if material.Article == "" {
		return entities.NewValidationError("article", "артикул обязателен")
	}
	if material.Name == "" {
		return entities.NewValidationError("name", "название обязательно")
	}
	if material.MaterialTypeID <= 0 {
		return entities.NewValidationError("material_type_id", "тип материала обязателен")
	}
	if material.MeasurementUnitID <= 0 {
		return entities.NewValidationError("measurement_unit_id", "единица измерения обязательна")
	}
	if material.PackageQuantity <= 0 {
		return entities.NewValidationError("package_quantity", "количество в упаковке должно быть больше нуля")
	}
	if material.CostPerUnit < 0 {
		return entities.NewValidationError("cost_per_unit", "стоимость за единицу не может быть отрицательной")
	}
	if material.StockQuantity < 0 {
		return entities.NewValidationError("stock_quantity", "остаток на складе не может быть отрицательным")
	}
	if material.MinStockQuantity < 0 {
		return entities.NewValidationError("min_stock_quantity", "минимальный остаток не может быть отрицательным")
	}
	return nil
}

// CalculatorUseCase содержит бизнес-логику для калькулятора материалов
type CalculatorUseCase struct {
	materialRepo repositories.MaterialRepository
}

// NewCalculatorUseCase создает новый use case калькулятора
func NewCalculatorUseCase(materialRepo repositories.MaterialRepository) *CalculatorUseCase {
	return &CalculatorUseCase{
		materialRepo: materialRepo,
	}
}

// CalculateRequiredMaterial рассчитывает необходимое количество материала
func (uc *CalculatorUseCase) CalculateRequiredMaterial(request *entities.MaterialCalculationRequest) (int, error) {
	// Валидация входных данных
	if err := uc.validateCalculationRequest(request); err != nil {
		return 0, fmt.Errorf("ошибка валидации: %w", err)
	}

	// Получаем тип продукции
	productType, err := uc.materialRepo.GetProductTypeByID(request.ProductTypeID)
	if err != nil {
		return 0, entities.NewNotFoundError("тип продукции", strconv.Itoa(request.ProductTypeID))
	}

	// Получаем тип материала
	materialType, err := uc.materialRepo.GetMaterialTypeByID(request.MaterialTypeID)
	if err != nil {
		return 0, entities.NewNotFoundError("тип материала", strconv.Itoa(request.MaterialTypeID))
	}

	// Используем доменную логику для расчета
	result := request.CalculateRequiredQuantity(
		productType.Coefficient,
		materialType.WastePercentage,
	)

	return result, nil
}

// validateCalculationRequest проверяет корректность запроса на расчет
func (uc *CalculatorUseCase) validateCalculationRequest(request *entities.MaterialCalculationRequest) error {
	if request.ProductTypeID <= 0 {
		return entities.NewValidationError("product_type_id", "ID типа продукции должен быть больше нуля")
	}
	if request.MaterialTypeID <= 0 {
		return entities.NewValidationError("material_type_id", "ID типа материала должен быть больше нуля")
	}
	if request.ProductQuantity <= 0 {
		return entities.NewValidationError("product_quantity", "количество продукции должно быть больше нуля")
	}
	if request.ProductParam1 <= 0 {
		return entities.NewValidationError("product_param1", "первый параметр продукции должен быть больше нуля")
	}
	if request.ProductParam2 <= 0 {
		return entities.NewValidationError("product_param2", "второй параметр продукции должен быть больше нуля")
	}
	if request.MaterialInStock < 0 {
		return entities.NewValidationError("material_in_stock", "количество материала на складе не может быть отрицательным")
	}
	return nil
}
