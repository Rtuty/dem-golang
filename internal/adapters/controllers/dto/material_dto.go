package dto

import "wallpaper-system/internal/domain/entities"

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
