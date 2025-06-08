package entity

import (
	"fmt"
	"time"
)

// ProductType представляет тип продукции
type ProductType struct {
	id          ID
	name        string
	description *string
	coefficient Money // Коэффициент для расчета материалов
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProductType создает новый тип продукции
func NewProductType(name string, coefficient Money) (*ProductType, error) {
	if err := validateProductTypeName(name); err != nil {
		return nil, err
	}

	if coefficient <= 0 {
		return nil, fmt.Errorf("коэффициент должен быть положительным числом")
	}

	return &ProductType{
		name:        name,
		coefficient: coefficient,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}, nil
}

// Getters
func (pt *ProductType) ID() ID                { return pt.id }
func (pt *ProductType) Name() string          { return pt.name }
func (pt *ProductType) Description() *string  { return pt.description }
func (pt *ProductType) Coefficient() Money    { return pt.coefficient }
func (pt *ProductType) CreatedAt() time.Time  { return pt.createdAt }
func (pt *ProductType) UpdatedAt() time.Time  { return pt.updatedAt }

// SetDescription устанавливает описание типа продукции
func (pt *ProductType) SetDescription(description *string) {
	pt.description = description
	pt.updatedAt = time.Now()
}

// SetCoefficient обновляет коэффициент
func (pt *ProductType) SetCoefficient(coefficient Money) error {
	if coefficient <= 0 {
		return fmt.Errorf("коэффициент должен быть положительным числом")
	}
	pt.coefficient = coefficient
	pt.updatedAt = time.Now()
	return nil
}

// validateProductTypeName проверяет корректность названия типа продукции
func validateProductTypeName(name string) error {
	if name == "" {
		return fmt.Errorf("название типа продукции не может быть пустым")
	}

	if len(name) > 100 {
		return fmt.Errorf("название типа продукции не может превышать 100 символов")
	}

	return nil
} 