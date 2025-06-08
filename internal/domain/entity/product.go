package entity

import (
	"errors"
	"time"
)

// Product представляет сущность продукции в доменном слое
type Product struct {
	id                     ID
	article                string
	productTypeID          ID
	name                   string
	description            *string
	imagePath              *string
	minPartnerPrice        Money
	packageDimensions      *PackageDimensions
	weights                *ProductWeights
	qualityCertificatePath *string
	standardNumber         *string
	productionTime         *ProductionTime
	costPrice              *Money
	workshopNumber         *string
	requiredWorkers        *int
	rollWidth              *float64
	createdAt              time.Time
	updatedAt              time.Time
	productType            *ProductType
	materials              []ProductMaterial
	calculatedPrice        *Money
}

// ID представляет идентификатор сущности
type ID int

// Money представляет денежную сумму
type Money float64

// PackageDimensions представляет размеры упаковки
type PackageDimensions struct {
	Length float64
	Width  float64
	Height float64
}

// ProductWeights представляет веса продукции
type ProductWeights struct {
	WithoutPackage float64
	WithPackage    float64
}

// ProductionTime представляет время производства
type ProductionTime struct {
	Hours int
}

// ProductMaterial представляет связь продукции с материалом
type ProductMaterial struct {
	id              ID
	productID       ID
	materialID      ID
	quantityPerUnit float64
	createdAt       time.Time
	material        *Material
}

// NewProduct создает новую сущность продукции
func NewProduct(
	article string,
	productTypeID ID,
	name string,
	minPartnerPrice Money,
) (*Product, error) {
	if err := validateProductData(article, name, minPartnerPrice); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Product{
		article:         article,
		productTypeID:   productTypeID,
		name:            name,
		minPartnerPrice: minPartnerPrice,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

// validateProductData проверяет корректность данных продукции
func validateProductData(article, name string, minPartnerPrice Money) error {
	if article == "" {
		return errors.New("артикул не может быть пустым")
	}
	if name == "" {
		return errors.New("наименование не может быть пустым")
	}
	if minPartnerPrice < 0 {
		return errors.New("минимальная цена для партнера не может быть отрицательной")
	}
	return nil
}

// Getters
func (p *Product) ID() ID                            { return p.id }
func (p *Product) Article() string                   { return p.article }
func (p *Product) ProductTypeID() ID                 { return p.productTypeID }
func (p *Product) Name() string                      { return p.name }
func (p *Product) Description() *string              { return p.description }
func (p *Product) ImagePath() *string                { return p.imagePath }
func (p *Product) MinPartnerPrice() Money            { return p.minPartnerPrice }
func (p *Product) PackageDimensions() *PackageDimensions { return p.packageDimensions }
func (p *Product) Weights() *ProductWeights          { return p.weights }
func (p *Product) QualityCertificatePath() *string   { return p.qualityCertificatePath }
func (p *Product) StandardNumber() *string           { return p.standardNumber }
func (p *Product) ProductionTime() *ProductionTime   { return p.productionTime }
func (p *Product) CostPrice() *Money                 { return p.costPrice }
func (p *Product) WorkshopNumber() *string           { return p.workshopNumber }
func (p *Product) RequiredWorkers() *int             { return p.requiredWorkers }
func (p *Product) RollWidth() *float64               { return p.rollWidth }
func (p *Product) CreatedAt() time.Time              { return p.createdAt }
func (p *Product) UpdatedAt() time.Time              { return p.updatedAt }
func (p *Product) ProductType() *ProductType         { return p.productType }
func (p *Product) Materials() []ProductMaterial      { return p.materials }
func (p *Product) CalculatedPrice() *Money           { return p.calculatedPrice }

// Business methods
func (p *Product) UpdateBasicInfo(name string, description *string) error {
	if name == "" {
		return errors.New("наименование не может быть пустым")
	}
	p.name = name
	p.description = description
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) UpdatePrice(minPartnerPrice Money) error {
	if minPartnerPrice < 0 {
		return errors.New("минимальная цена для партнера не может быть отрицательной")
	}
	p.minPartnerPrice = minPartnerPrice
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) SetPackageDimensions(length, width, height float64) error {
	if length <= 0 || width <= 0 || height <= 0 {
		return errors.New("размеры упаковки должны быть положительными")
	}
	p.packageDimensions = &PackageDimensions{
		Length: length,
		Width:  width,
		Height: height,
	}
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) SetWeights(withoutPackage, withPackage float64) error {
	if withoutPackage <= 0 || withPackage <= 0 {
		return errors.New("веса должны быть положительными")
	}
	if withPackage <= withoutPackage {
		return errors.New("вес с упаковкой должен быть больше веса без упаковки")
	}
	p.weights = &ProductWeights{
		WithoutPackage: withoutPackage,
		WithPackage:    withPackage,
	}
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) SetProductionTime(hours int) error {
	if hours <= 0 {
		return errors.New("время производства должно быть положительным")
	}
	p.productionTime = &ProductionTime{Hours: hours}
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) CalculatePrice() error {
	if len(p.materials) == 0 {
		return errors.New("невозможно рассчитать цену: отсутствуют материалы")
	}

	var totalCost float64
	for _, material := range p.materials {
		if material.material == nil {
			continue
		}
		materialCost := material.quantityPerUnit * float64(material.material.CostPerUnit())
		totalCost += materialCost
	}

	calculatedPrice := Money(totalCost)
	p.calculatedPrice = &calculatedPrice
	return nil
}

// SetID устанавливает ID (используется репозиторием после сохранения)
func (p *Product) SetID(id ID) {
	p.id = id
}

// SetProductType устанавливает тип продукции
func (p *Product) SetProductType(productType *ProductType) {
	p.productType = productType
}

// SetMaterials устанавливает материалы
func (p *Product) SetMaterials(materials []ProductMaterial) {
	p.materials = materials
}

// Методы для ProductMaterial
func (pm *ProductMaterial) ID() ID              { return pm.id }
func (pm *ProductMaterial) ProductID() ID       { return pm.productID }
func (pm *ProductMaterial) MaterialID() ID      { return pm.materialID }
func (pm *ProductMaterial) QuantityPerUnit() float64 { return pm.quantityPerUnit }
func (pm *ProductMaterial) CreatedAt() time.Time { return pm.createdAt }
func (pm *ProductMaterial) Material() *Material { return pm.material }

func (pm *ProductMaterial) SetMaterial(material *Material) {
	pm.material = material
}

// Конструктор для ProductMaterial
func NewProductMaterial(productID, materialID ID, quantityPerUnit float64) (*ProductMaterial, error) {
	if quantityPerUnit <= 0 {
		return nil, errors.New("количество материала должно быть положительным")
	}

	return &ProductMaterial{
		productID:       productID,
		materialID:      materialID,
		quantityPerUnit: quantityPerUnit,
		createdAt:       time.Now(),
	}, nil
} 