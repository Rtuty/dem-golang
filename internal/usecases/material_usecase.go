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
		return -1, fmt.Errorf("ошибка валидации: %w", err)
	}

	// Получаем тип продукции
	productType, err := uc.materialRepo.GetProductTypeByID(request.ProductTypeID)
	if err != nil {
		return -1, entities.NewNotFoundError("тип продукции", strconv.Itoa(request.ProductTypeID))
	}

	// Получаем тип материала
	materialType, err := uc.materialRepo.GetMaterialTypeByID(request.MaterialTypeID)
	if err != nil {
		return -1, entities.NewNotFoundError("тип материала", strconv.Itoa(request.MaterialTypeID))
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
