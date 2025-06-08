package entity

import (
	"errors"
	"time"
)

// Material представляет сущность материала в доменном слое
type Material struct {
	id                ID
	article           string
	materialTypeID    ID
	name              string
	description       *string
	measurementUnitID ID
	packageQuantity   float64
	costPerUnit       Money
	stockQuantity     float64
	minStockQuantity  float64
	imagePath         *string
	createdAt         time.Time
	updatedAt         time.Time
	materialType      *MaterialType
	measurementUnit   *MeasurementUnit
}

// MaterialType представляет тип материала
type MaterialType struct {
	id              ID
	name            string
	wastePercentage float64
	createdAt       time.Time
	updatedAt       time.Time
}

// MeasurementUnit представляет единицу измерения
type MeasurementUnit struct {
	id           ID
	name         string
	abbreviation string
	createdAt    time.Time
}

// ProductType представляет тип продукции
type ProductType struct {
	id          ID
	name        string
	coefficient float64
	createdAt   time.Time
	updatedAt   time.Time
}

// NewMaterial создает новую сущность материала
func NewMaterial(
	article string,
	materialTypeID ID,
	name string,
	measurementUnitID ID,
	packageQuantity float64,
	costPerUnit Money,
) (*Material, error) {
	if err := validateMaterialData(article, name, packageQuantity, costPerUnit); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Material{
		article:           article,
		materialTypeID:    materialTypeID,
		name:              name,
		measurementUnitID: measurementUnitID,
		packageQuantity:   packageQuantity,
		costPerUnit:       costPerUnit,
		stockQuantity:     0,
		minStockQuantity:  0,
		createdAt:         now,
		updatedAt:         now,
	}, nil
}

// validateMaterialData проверяет корректность данных материала
func validateMaterialData(article, name string, packageQuantity float64, costPerUnit Money) error {
	if article == "" {
		return errors.New("артикул материала не может быть пустым")
	}
	if name == "" {
		return errors.New("наименование материала не может быть пустым")
	}
	if packageQuantity <= 0 {
		return errors.New("количество в упаковке должно быть положительным")
	}
	if costPerUnit < 0 {
		return errors.New("стоимость за единицу не может быть отрицательной")
	}
	return nil
}

// Material getters
func (m *Material) ID() ID                      { return m.id }
func (m *Material) Article() string             { return m.article }
func (m *Material) MaterialTypeID() ID          { return m.materialTypeID }
func (m *Material) Name() string                { return m.name }
func (m *Material) Description() *string        { return m.description }
func (m *Material) MeasurementUnitID() ID       { return m.measurementUnitID }
func (m *Material) PackageQuantity() float64    { return m.packageQuantity }
func (m *Material) CostPerUnit() Money          { return m.costPerUnit }
func (m *Material) StockQuantity() float64      { return m.stockQuantity }
func (m *Material) MinStockQuantity() float64   { return m.minStockQuantity }
func (m *Material) ImagePath() *string          { return m.imagePath }
func (m *Material) CreatedAt() time.Time        { return m.createdAt }
func (m *Material) UpdatedAt() time.Time        { return m.updatedAt }
func (m *Material) MaterialType() *MaterialType { return m.materialType }
func (m *Material) MeasurementUnit() *MeasurementUnit { return m.measurementUnit }

// Business methods for Material
func (m *Material) UpdateBasicInfo(name string, description *string) error {
	if name == "" {
		return errors.New("наименование материала не может быть пустым")
	}
	m.name = name
	m.description = description
	m.updatedAt = time.Now()
	return nil
}

func (m *Material) UpdateCost(costPerUnit Money) error {
	if costPerUnit < 0 {
		return errors.New("стоимость за единицу не может быть отрицательной")
	}
	m.costPerUnit = costPerUnit
	m.updatedAt = time.Now()
	return nil
}

func (m *Material) UpdateStock(newStockQuantity float64) error {
	if newStockQuantity < 0 {
		return errors.New("количество на складе не может быть отрицательным")
	}
	m.stockQuantity = newStockQuantity
	m.updatedAt = time.Now()
	return nil
}

func (m *Material) AddToStock(quantity float64) error {
	if quantity <= 0 {
		return errors.New("количество для добавления должно быть положительным")
	}
	m.stockQuantity += quantity
	m.updatedAt = time.Now()
	return nil
}

func (m *Material) RemoveFromStock(quantity float64) error {
	if quantity <= 0 {
		return errors.New("количество для списания должно быть положительным")
	}
	if m.stockQuantity < quantity {
		return errors.New("недостаточно материала на складе")
	}
	m.stockQuantity -= quantity
	m.updatedAt = time.Now()
	return nil
}

func (m *Material) SetMinStockQuantity(minQuantity float64) error {
	if minQuantity < 0 {
		return errors.New("минимальное количество на складе не может быть отрицательным")
	}
	m.minStockQuantity = minQuantity
	m.updatedAt = time.Now()
	return nil
}

func (m *Material) IsLowStock() bool {
	return m.stockQuantity <= m.minStockQuantity
}

// SetID устанавливает ID (используется репозиторием после сохранения)
func (m *Material) SetID(id ID) {
	m.id = id
}

// SetMaterialType устанавливает тип материала
func (m *Material) SetMaterialType(materialType *MaterialType) {
	m.materialType = materialType
}

// SetMeasurementUnit устанавливает единицу измерения
func (m *Material) SetMeasurementUnit(unit *MeasurementUnit) {
	m.measurementUnit = unit
}

// MaterialType methods
func NewMaterialType(name string, wastePercentage float64) (*MaterialType, error) {
	if name == "" {
		return nil, errors.New("наименование типа материала не может быть пустым")
	}
	if wastePercentage < 0 || wastePercentage > 100 {
		return nil, errors.New("процент брака должен быть от 0 до 100")
	}

	now := time.Now()
	return &MaterialType{
		name:            name,
		wastePercentage: wastePercentage,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func (mt *MaterialType) ID() ID                { return mt.id }
func (mt *MaterialType) Name() string          { return mt.name }
func (mt *MaterialType) WastePercentage() float64 { return mt.wastePercentage }
func (mt *MaterialType) CreatedAt() time.Time  { return mt.createdAt }
func (mt *MaterialType) UpdatedAt() time.Time  { return mt.updatedAt }

func (mt *MaterialType) SetID(id ID) {
	mt.id = id
}

func (mt *MaterialType) UpdateWastePercentage(percentage float64) error {
	if percentage < 0 || percentage > 100 {
		return errors.New("процент брака должен быть от 0 до 100")
	}
	mt.wastePercentage = percentage
	mt.updatedAt = time.Now()
	return nil
}

// MeasurementUnit methods
func NewMeasurementUnit(name, abbreviation string) (*MeasurementUnit, error) {
	if name == "" {
		return nil, errors.New("наименование единицы измерения не может быть пустым")
	}
	if abbreviation == "" {
		return nil, errors.New("сокращение единицы измерения не может быть пустым")
	}

	return &MeasurementUnit{
		name:         name,
		abbreviation: abbreviation,
		createdAt:    time.Now(),
	}, nil
}

func (mu *MeasurementUnit) ID() ID             { return mu.id }
func (mu *MeasurementUnit) Name() string       { return mu.name }
func (mu *MeasurementUnit) Abbreviation() string { return mu.abbreviation }
func (mu *MeasurementUnit) CreatedAt() time.Time { return mu.createdAt }

func (mu *MeasurementUnit) SetID(id ID) {
	mu.id = id
}

// ProductType methods
func NewProductType(name string, coefficient float64) (*ProductType, error) {
	if name == "" {
		return nil, errors.New("наименование типа продукции не может быть пустым")
	}
	if coefficient <= 0 {
		return nil, errors.New("коэффициент должен быть положительным")
	}

	now := time.Now()
	return &ProductType{
		name:        name,
		coefficient: coefficient,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (pt *ProductType) ID() ID              { return pt.id }
func (pt *ProductType) Name() string        { return pt.name }
func (pt *ProductType) Coefficient() float64 { return pt.coefficient }
func (pt *ProductType) CreatedAt() time.Time { return pt.createdAt }
func (pt *ProductType) UpdatedAt() time.Time { return pt.updatedAt }

func (pt *ProductType) SetID(id ID) {
	pt.id = id
}

func (pt *ProductType) UpdateCoefficient(coefficient float64) error {
	if coefficient <= 0 {
		return errors.New("коэффициент должен быть положительным")
	}
	pt.coefficient = coefficient
	pt.updatedAt = time.Now()
	return nil
} 