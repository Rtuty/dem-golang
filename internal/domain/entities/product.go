package entities

import (
	"time"
)

// ProductType представляет тип продукции в предметной области
type ProductType struct {
	ID          int
	Name        string
	Coefficient float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Product представляет продукцию в предметной области
type Product struct {
	ID                     int
	Article                string
	ProductTypeID          int
	Name                   string
	Description            *string
	ImagePath              *string
	MinPartnerPrice        float64
	PackageLength          *float64
	PackageWidth           *float64
	PackageHeight          *float64
	WeightWithoutPackage   *float64
	WeightWithPackage      *float64
	QualityCertificatePath *string
	StandardNumber         *string
	ProductionTimeHours    *float64
	CostPrice              *float64
	WorkshopNumber         *string
	RequiredWorkers        *int
	RollWidth              *float64
	CreatedAt              time.Time
	UpdatedAt              time.Time

	// Связанные данные
	ProductType     *ProductType
	Materials       []ProductMaterial
	CalculatedPrice *float64
}

// ProductMaterial представляет связь продукции с материалами
type ProductMaterial struct {
	ID              int
	ProductID       int
	MaterialID      int
	QuantityPerUnit float64
	CreatedAt       time.Time
	Material        *Material
}

// CalculatePrice рассчитывает стоимость продукции на основе материалов
func (p *Product) CalculatePrice() float64 {
	if len(p.Materials) == 0 || p.ProductType == nil {
		return 0
	}

	materialsCost := 0.0
	for _, pm := range p.Materials {
		if pm.Material != nil {
			materialsCost += pm.QuantityPerUnit * pm.Material.CostPerUnit
		}
	}

	// Применяем коэффициент типа продукции
	materialsCost *= p.ProductType.Coefficient

	// Добавляем 20% наценку
	return materialsCost * 1.2
}

// Validate проверяет корректность данных продукции
func (p *Product) Validate() error {
	if p.Article == "" {
		return NewValidationError("article", "артикул не может быть пустым")
	}
	if p.Name == "" {
		return NewValidationError("name", "название не может быть пустым")
	}
	if p.MinPartnerPrice < 0 {
		return NewValidationError("min_partner_price", "минимальная партнерская цена не может быть отрицательной")
	}
	if p.RollWidth != nil && *p.RollWidth < 0 {
		return NewValidationError("roll_width", "ширина рулона не может быть отрицательной")
	}
	return nil
}
