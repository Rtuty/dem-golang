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
