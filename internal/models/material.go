package models

import (
	"time"
)

// MaterialType представляет типы материалов
type MaterialType struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	WastePercentage float64   `json:"waste_percentage" db:"waste_percentage"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// MeasurementUnit представляет единицы измерения
type MeasurementUnit struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Abbreviation string    `json:"abbreviation" db:"abbreviation"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Material представляет материалы (сырье)
type Material struct {
	ID                int       `json:"id" db:"id"`
	Article           string    `json:"article" db:"article"`
	MaterialTypeID    int       `json:"material_type_id" db:"material_type_id"`
	Name              string    `json:"name" db:"name"`
	Description       *string   `json:"description" db:"description"`
	MeasurementUnitID int       `json:"measurement_unit_id" db:"measurement_unit_id"`
	PackageQuantity   float64   `json:"package_quantity" db:"package_quantity"`
	CostPerUnit       float64   `json:"cost_per_unit" db:"cost_per_unit"`
	StockQuantity     float64   `json:"stock_quantity" db:"stock_quantity"`
	MinStockQuantity  float64   `json:"min_stock_quantity" db:"min_stock_quantity"`
	ImagePath         *string   `json:"image_path" db:"image_path"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`

	// Связанные данные
	MaterialType    *MaterialType    `json:"material_type,omitempty"`
	MeasurementUnit *MeasurementUnit `json:"measurement_unit,omitempty"`
}

// MaterialForProduct представляет материал в контексте продукции
type MaterialForProduct struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Article          string  `json:"article"`
	QuantityPerUnit  float64 `json:"quantity_per_unit"`
	CostPerUnit      float64 `json:"cost_per_unit"`
	TotalCost        float64 `json:"total_cost"`
	UnitAbbreviation string  `json:"unit_abbreviation"`
}

// MaterialStockHistory представляет историю движения материалов на складе
type MaterialStockHistory struct {
	ID                int       `json:"id" db:"id"`
	MaterialID        int       `json:"material_id" db:"material_id"`
	OperationType     string    `json:"operation_type" db:"operation_type"`
	Quantity          float64   `json:"quantity" db:"quantity"`
	OldStock          float64   `json:"old_stock" db:"old_stock"`
	NewStock          float64   `json:"new_stock" db:"new_stock"`
	ReferenceDocument *string   `json:"reference_document" db:"reference_document"`
	Comment           *string   `json:"comment" db:"comment"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`

	// Связанные данные
	Material *Material `json:"material,omitempty"`
}

// MaterialCalculationRequest представляет запрос на расчет необходимого количества материала
type MaterialCalculationRequest struct {
	ProductTypeID   int     `json:"product_type_id" binding:"required"`
	MaterialTypeID  int     `json:"material_type_id" binding:"required"`
	ProductQuantity int     `json:"product_quantity" binding:"required,min=1"`
	ProductParam1   float64 `json:"product_param1" binding:"required,min=0"`
	ProductParam2   float64 `json:"product_param2" binding:"required,min=0"`
	MaterialInStock float64 `json:"material_in_stock" binding:"min=0"`
}

// MaterialCalculationResponse представляет ответ на расчет материала
type MaterialCalculationResponse struct {
	RequiredQuantity int `json:"required_quantity"`
}
