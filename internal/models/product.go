package models

import (
	"time"
)

// ProductType представляет типы продукции
type ProductType struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Coefficient float64   `json:"coefficient" db:"coefficient"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Product представляет продукцию компании
type Product struct {
	ID                     int       `json:"id" db:"id"`
	Article                string    `json:"article" db:"article"`
	ProductTypeID          int       `json:"product_type_id" db:"product_type_id"`
	Name                   string    `json:"name" db:"name"`
	Description            *string   `json:"description" db:"description"`
	ImagePath              *string   `json:"image_path" db:"image_path"`
	MinPartnerPrice        float64   `json:"min_partner_price" db:"min_partner_price"`
	PackageLength          *float64  `json:"package_length" db:"package_length"`
	PackageWidth           *float64  `json:"package_width" db:"package_width"`
	PackageHeight          *float64  `json:"package_height" db:"package_height"`
	WeightWithoutPackage   *float64  `json:"weight_without_package" db:"weight_without_package"`
	WeightWithPackage      *float64  `json:"weight_with_package" db:"weight_with_package"`
	QualityCertificatePath *string   `json:"quality_certificate_path" db:"quality_certificate_path"`
	StandardNumber         *string   `json:"standard_number" db:"standard_number"`
	ProductionTimeHours    *int      `json:"production_time_hours" db:"production_time_hours"`
	CostPrice              *float64  `json:"cost_price" db:"cost_price"`
	WorkshopNumber         *string   `json:"workshop_number" db:"workshop_number"`
	RequiredWorkers        *int      `json:"required_workers" db:"required_workers"`
	RollWidth              *float64  `json:"roll_width" db:"roll_width"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`

	// Связанные данные
	ProductType     *ProductType      `json:"product_type,omitempty"`
	Materials       []ProductMaterial `json:"materials,omitempty"`
	CalculatedPrice *float64          `json:"calculated_price,omitempty"`
}

// ProductMaterial представляет связь продукции с материалами
type ProductMaterial struct {
	ID              int       `json:"id" db:"id"`
	ProductID       int       `json:"product_id" db:"product_id"`
	MaterialID      int       `json:"material_id" db:"material_id"`
	QuantityPerUnit float64   `json:"quantity_per_unit" db:"quantity_per_unit"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`

	// Связанные данные
	Material *Material `json:"material,omitempty"`
}

// ProductWithType представляет продукцию с информацией о типе
type ProductWithType struct {
	Product
	TypeName        string  `json:"type_name" db:"type_name"`
	TypeCoefficient float64 `json:"type_coefficient" db:"type_coefficient"`
}

// ProductListItem представляет элемент списка продукции для отображения
type ProductListItem struct {
	ID              int      `json:"id"`
	Article         string   `json:"article"`
	TypeName        string   `json:"type_name"`
	Name            string   `json:"name"`
	MinPartnerPrice float64  `json:"min_partner_price"`
	RollWidth       *float64 `json:"roll_width"`
	CalculatedPrice *float64 `json:"calculated_price"`
}

// CreateProductRequest представляет запрос на создание продукции
type CreateProductRequest struct {
	Article         string   `json:"article" binding:"required"`
	ProductTypeID   int      `json:"product_type_id" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Description     *string  `json:"description"`
	MinPartnerPrice float64  `json:"min_partner_price" binding:"required,min=0"`
	RollWidth       *float64 `json:"roll_width" binding:"omitempty,min=0"`
}

// UpdateProductRequest представляет запрос на обновление продукции
type UpdateProductRequest struct {
	Article         *string  `json:"article"`
	ProductTypeID   *int     `json:"product_type_id"`
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	MinPartnerPrice *float64 `json:"min_partner_price" binding:"omitempty,min=0"`
	RollWidth       *float64 `json:"roll_width" binding:"omitempty,min=0"`
}
