package entities

import (
	"time"
)

// MaterialType представляет тип материала в предметной области
type MaterialType struct {
	ID              int
	Name            string
	WastePercentage float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// MeasurementUnit представляет единицу измерения
type MeasurementUnit struct {
	ID           int
	Name         string
	Abbreviation string
	CreatedAt    time.Time
}

// Material представляет материал в предметной области
type Material struct {
	ID                int
	Article           string
	MaterialTypeID    int
	Name              string
	Description       *string
	MeasurementUnitID int
	PackageQuantity   float64
	CostPerUnit       float64
	StockQuantity     float64
	MinStockQuantity  float64
	ImagePath         *string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Связанные данные
	MaterialType    *MaterialType
	MeasurementUnit *MeasurementUnit
}

// MaterialCalculationRequest представляет запрос на расчет материала
type MaterialCalculationRequest struct {
	ProductTypeID   int
	MaterialTypeID  int
	ProductQuantity int
	ProductParam1   float64
	ProductParam2   float64
	MaterialInStock float64
}

// CalculateRequiredQuantity рассчитывает необходимое количество материала
func (mcr *MaterialCalculationRequest) CalculateRequiredQuantity(
	productTypeCoeff float64,
	wastePercentage float64,
) int {
	if mcr.ProductQuantity <= 0 || mcr.ProductParam1 <= 0 || mcr.ProductParam2 <= 0 {
		return -1
	}

	// Рассчитываем материал на единицу продукции
	materialPerUnit := mcr.ProductParam1 * mcr.ProductParam2 * productTypeCoeff

	// Общее количество материала
	totalMaterial := materialPerUnit * float64(mcr.ProductQuantity)

	// Учитываем процент брака
	materialWithWaste := totalMaterial * (1 + wastePercentage/100)

	// Учитываем остатки на складе
	requiredMaterial := materialWithWaste - mcr.MaterialInStock

	if requiredMaterial <= 0 {
		return 0
	}

	// Округляем в большую сторону до целого числа
	return int(requiredMaterial + 0.9999999)
}

// Validate проверяет корректность данных материала
func (m *Material) Validate() error {
	if m.Article == "" {
		return NewValidationError("article", "артикул не может быть пустым")
	}
	if m.Name == "" {
		return NewValidationError("name", "название не может быть пустым")
	}
	if m.CostPerUnit < 0 {
		return NewValidationError("cost_per_unit", "стоимость не может быть отрицательной")
	}
	if m.PackageQuantity <= 0 {
		return NewValidationError("package_quantity", "количество в упаковке должно быть больше нуля")
	}
	if m.StockQuantity < 0 {
		return NewValidationError("stock_quantity", "количество на складе не может быть отрицательным")
	}
	return nil
}

// CalculateRequiredQuantity рассчитывает необходимое количество материала с учетом отходов
func (m *Material) CalculateRequiredQuantity(baseQuantity, wastePercentage float64) (int, error) {
	if baseQuantity < 0 {
		return 0, NewValidationError("base_quantity", "Базовое количество не может быть отрицательным")
	}

	if wastePercentage < 0 {
		return 0, NewValidationError("waste_percentage", "Процент отходов не может быть отрицательным")
	}

	if baseQuantity == 0 {
		return 0, nil
	}

	// Рассчитываем с учетом отходов
	requiredQuantity := baseQuantity * (1 + wastePercentage/100)

	// Округляем в большую сторону
	return int(requiredQuantity + 0.9999999), nil
}

// Validate проверяет валидность запроса на расчет материала
func (r *MaterialCalculationRequest) Validate() error {
	if r.ProductTypeID <= 0 {
		return NewValidationError("product_type_id", "ID типа продукции должен быть больше нуля")
	}

	if r.MaterialTypeID <= 0 {
		return NewValidationError("material_type_id", "ID типа материала должен быть больше нуля")
	}

	if r.ProductQuantity <= 0 {
		return NewValidationError("product_quantity", "Количество продукции должно быть больше нуля")
	}

	if r.ProductParam1 < 0 {
		return NewValidationError("product_param1", "Первый параметр продукции не может быть отрицательным")
	}

	if r.ProductParam2 < 0 {
		return NewValidationError("product_param2", "Второй параметр продукции не может быть отрицательным")
	}

	if r.MaterialInStock < 0 {
		return NewValidationError("material_in_stock", "Остаток материала на складе не может быть отрицательным")
	}

	return nil
}
