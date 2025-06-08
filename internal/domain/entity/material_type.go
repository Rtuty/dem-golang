package entity

import (
	"fmt"
	"time"
)

// MaterialType представляет тип материала
type MaterialType struct {
	id         ID
	name       string
	description *string
	defectRate float64 // Процент брака (0.05 = 5%)
	createdAt  time.Time
	updatedAt  time.Time
}

// NewMaterialType создает новый тип материала
func NewMaterialType(name string, defectRate float64) (*MaterialType, error) {
	if err := validateMaterialTypeName(name); err != nil {
		return nil, err
	}

	if err := validateDefectRate(defectRate); err != nil {
		return nil, err
	}

	return &MaterialType{
		name:       name,
		defectRate: defectRate,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}, nil
}

// Getters
func (mt *MaterialType) ID() ID                { return mt.id }
func (mt *MaterialType) Name() string          { return mt.name }
func (mt *MaterialType) Description() *string  { return mt.description }
func (mt *MaterialType) DefectRate() float64   { return mt.defectRate }
func (mt *MaterialType) CreatedAt() time.Time  { return mt.createdAt }
func (mt *MaterialType) UpdatedAt() time.Time  { return mt.updatedAt }

// SetDescription устанавливает описание типа материала
func (mt *MaterialType) SetDescription(description *string) {
	mt.description = description
	mt.updatedAt = time.Now()
}

// SetDefectRate обновляет процент брака
func (mt *MaterialType) SetDefectRate(defectRate float64) error {
	if err := validateDefectRate(defectRate); err != nil {
		return err
	}
	mt.defectRate = defectRate
	mt.updatedAt = time.Now()
	return nil
}

// CalculateWithDefect рассчитывает количество с учетом брака
func (mt *MaterialType) CalculateWithDefect(baseQuantity float64) float64 {
	return baseQuantity * (1.0 + mt.defectRate)
}

// validateMaterialTypeName проверяет корректность названия типа материала
func validateMaterialTypeName(name string) error {
	if name == "" {
		return fmt.Errorf("название типа материала не может быть пустым")
	}

	if len(name) > 100 {
		return fmt.Errorf("название типа материала не может превышать 100 символов")
	}

	return nil
}

// validateDefectRate проверяет корректность процента брака
func validateDefectRate(defectRate float64) error {
	if defectRate < 0 {
		return fmt.Errorf("процент брака не может быть отрицательным")
	}

	if defectRate > 1.0 {
		return fmt.Errorf("процент брака не может превышать 100%%")
	}

	return nil
} 