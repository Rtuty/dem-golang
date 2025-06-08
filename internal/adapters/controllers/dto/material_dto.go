package dto

import (
	"wallpaper-system/internal/domain/entities"
)

// MaterialCalculationRequestDTO представляет запрос на расчет материала через API
type MaterialCalculationRequestDTO struct {
	ProductTypeID   int     `json:"product_type_id" binding:"required"`
	MaterialTypeID  int     `json:"material_type_id" binding:"required"`
	ProductQuantity int     `json:"product_quantity" binding:"required,min=1"`
	ProductParam1   float64 `json:"product_param1" binding:"required,min=0"`
	ProductParam2   float64 `json:"product_param2" binding:"required,min=0"`
	MaterialInStock float64 `json:"material_in_stock" binding:"min=0"`
}

// MaterialCalculationResponseDTO представляет ответ на расчет материала
type MaterialCalculationResponseDTO struct {
	RequiredQuantity int `json:"required_quantity"`
}

// ToEntity преобразует DTO в доменную сущность
func (dto *MaterialCalculationRequestDTO) ToEntity() *entities.MaterialCalculationRequest {
	return &entities.MaterialCalculationRequest{
		ProductTypeID:   dto.ProductTypeID,
		MaterialTypeID:  dto.MaterialTypeID,
		ProductQuantity: dto.ProductQuantity,
		ProductParam1:   dto.ProductParam1,
		ProductParam2:   dto.ProductParam2,
		MaterialInStock: dto.MaterialInStock,
	}
}

// CreateMaterialDTO представляет данные для создания материала
type CreateMaterialDTO struct {
	Article             string  `form:"article" json:"article" binding:"required"`
	MaterialTypeID      int     `form:"material_type_id" json:"material_type_id" binding:"required"`
	Name                string  `form:"name" json:"name" binding:"required"`
	Description         string  `form:"description" json:"description"`
	MeasurementUnitID   int     `form:"measurement_unit_id" json:"measurement_unit_id" binding:"required"`
	PackageQuantity     float64 `form:"package_quantity" json:"package_quantity" binding:"required,min=0"`
	CostPerUnit         float64 `form:"cost_per_unit" json:"cost_per_unit" binding:"required,min=0"`
	StockQuantity       float64 `form:"stock_quantity" json:"stock_quantity" binding:"min=0"`
	MinStockQuantity    float64 `form:"min_stock_quantity" json:"min_stock_quantity" binding:"min=0"`
	ImagePath           string  `form:"image_path" json:"image_path"`
}

// UpdateMaterialDTO представляет данные для обновления материала
type UpdateMaterialDTO struct {
	Article             string  `form:"article" json:"article" binding:"required"`
	MaterialTypeID      int     `form:"material_type_id" json:"material_type_id" binding:"required"`
	Name                string  `form:"name" json:"name" binding:"required"`
	Description         string  `form:"description" json:"description"`
	MeasurementUnitID   int     `form:"measurement_unit_id" json:"measurement_unit_id" binding:"required"`
	PackageQuantity     float64 `form:"package_quantity" json:"package_quantity" binding:"required,min=0"`
	CostPerUnit         float64 `form:"cost_per_unit" json:"cost_per_unit" binding:"required,min=0"`
	StockQuantity       float64 `form:"stock_quantity" json:"stock_quantity" binding:"min=0"`
	MinStockQuantity    float64 `form:"min_stock_quantity" json:"min_stock_quantity" binding:"min=0"`
	ImagePath           string  `form:"image_path" json:"image_path"`
}

// ToEntity преобразует CreateMaterialDTO в доменную сущность
func (dto *CreateMaterialDTO) ToEntity() *entities.Material {
	var description *string
	if dto.Description != "" {
		description = &dto.Description
	}
	
	var imagePath *string
	if dto.ImagePath != "" {
		imagePath = &dto.ImagePath
	}

	return &entities.Material{
		Article:             dto.Article,
		MaterialTypeID:      dto.MaterialTypeID,
		Name:                dto.Name,
		Description:         description,
		MeasurementUnitID:   dto.MeasurementUnitID,
		PackageQuantity:     dto.PackageQuantity,
		CostPerUnit:         dto.CostPerUnit,
		StockQuantity:       dto.StockQuantity,
		MinStockQuantity:    dto.MinStockQuantity,
		ImagePath:           imagePath,
	}
}

// ToEntity преобразует UpdateMaterialDTO в доменную сущность
func (dto *UpdateMaterialDTO) ToEntity() *entities.Material {
	var description *string
	if dto.Description != "" {
		description = &dto.Description
	}
	
	var imagePath *string
	if dto.ImagePath != "" {
		imagePath = &dto.ImagePath
	}

	return &entities.Material{
		Article:             dto.Article,
		MaterialTypeID:      dto.MaterialTypeID,
		Name:                dto.Name,
		Description:         description,
		MeasurementUnitID:   dto.MeasurementUnitID,
		PackageQuantity:     dto.PackageQuantity,
		CostPerUnit:         dto.CostPerUnit,
		StockQuantity:       dto.StockQuantity,
		MinStockQuantity:    dto.MinStockQuantity,
		ImagePath:           imagePath,
	}
}
