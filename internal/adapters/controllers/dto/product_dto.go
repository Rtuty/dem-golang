package dto

import (
	"wallpaper-system/internal/domain/entities"
)

// ProductListItemDTO представляет элемент списка продукции для API
type ProductListItemDTO struct {
	ID              int      `json:"id"`
	Article         string   `json:"article"`
	TypeName        string   `json:"type_name"`
	Name            string   `json:"name"`
	MinPartnerPrice float64  `json:"min_partner_price"`
	RollWidth       *float64 `json:"roll_width"`
	CalculatedPrice *float64 `json:"calculated_price"`
}

// ProductDetailDTO представляет детальную информацию о продукции
type ProductDetailDTO struct {
	ProductListItemDTO
	Description            *string              `json:"description"`
	ImagePath              *string              `json:"image_path"`
	PackageLength          *float64             `json:"package_length"`
	PackageWidth           *float64             `json:"package_width"`
	PackageHeight          *float64             `json:"package_height"`
	WeightWithoutPackage   *float64             `json:"weight_without_package"`
	WeightWithPackage      *float64             `json:"weight_with_package"`
	QualityCertificatePath *string              `json:"quality_certificate_path"`
	StandardNumber         *string              `json:"standard_number"`
	ProductionTimeHours    *float64             `json:"production_time_hours"`
	CostPrice              *float64             `json:"cost_price"`
	WorkshopNumber         *string              `json:"workshop_number"`
	RequiredWorkers        *int                 `json:"required_workers"`
	Materials              []ProductMaterialDTO `json:"materials"`
}

// ProductMaterialDTO представляет материал в контексте продукции
type ProductMaterialDTO struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Article          string  `json:"article"`
	QuantityPerUnit  float64 `json:"quantity_per_unit"`
	CostPerUnit      float64 `json:"cost_per_unit"`
	TotalCost        float64 `json:"total_cost"`
	UnitAbbreviation string  `json:"unit_abbreviation"`
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

// FromProductEntity преобразует доменную сущность в DTO для списка
func FromProductEntity(product *entities.Product) ProductListItemDTO {
	dto := ProductListItemDTO{
		ID:              product.ID,
		Article:         product.Article,
		Name:            product.Name,
		MinPartnerPrice: product.MinPartnerPrice,
		RollWidth:       product.RollWidth,
		CalculatedPrice: product.CalculatedPrice,
	}

	if product.ProductType != nil {
		dto.TypeName = product.ProductType.Name
	}

	return dto
}

// FromProductEntityWithMaterials преобразует доменную сущность в детальный DTO
func FromProductEntityWithMaterials(product *entities.Product, materials []entities.Material) ProductDetailDTO {
	dto := ProductDetailDTO{
		ProductListItemDTO:     FromProductEntity(product),
		Description:            product.Description,
		ImagePath:              product.ImagePath,
		PackageLength:          product.PackageLength,
		PackageWidth:           product.PackageWidth,
		PackageHeight:          product.PackageHeight,
		WeightWithoutPackage:   product.WeightWithoutPackage,
		WeightWithPackage:      product.WeightWithPackage,
		QualityCertificatePath: product.QualityCertificatePath,
		StandardNumber:         product.StandardNumber,
		ProductionTimeHours:    product.ProductionTimeHours,
		CostPrice:              product.CostPrice,
		WorkshopNumber:         product.WorkshopNumber,
		RequiredWorkers:        product.RequiredWorkers,
	}

	// Преобразуем материалы
	dto.Materials = make([]ProductMaterialDTO, len(product.Materials))
	for i, pm := range product.Materials {
		dto.Materials[i] = ProductMaterialDTO{
			ID:              pm.Material.ID,
			Name:            pm.Material.Name,
			Article:         pm.Material.Article,
			QuantityPerUnit: pm.QuantityPerUnit,
			CostPerUnit:     pm.Material.CostPerUnit,
			TotalCost:       pm.QuantityPerUnit * pm.Material.CostPerUnit,
		}
		if pm.Material.MeasurementUnit != nil {
			dto.Materials[i].UnitAbbreviation = pm.Material.MeasurementUnit.Abbreviation
		}
	}

	return dto
}

// ToEntity преобразует DTO в доменную сущность
func (dto *CreateProductRequest) ToEntity() *entities.Product {
	return &entities.Product{
		Article:         dto.Article,
		ProductTypeID:   dto.ProductTypeID,
		Name:            dto.Name,
		Description:     dto.Description,
		MinPartnerPrice: dto.MinPartnerPrice,
		RollWidth:       dto.RollWidth,
	}
}

// ToEntity преобразует DTO в доменную сущность
func (dto *UpdateProductRequest) ToEntity(id int) *entities.Product {
	product := &entities.Product{ID: id}

	if dto.Article != nil {
		product.Article = *dto.Article
	}
	if dto.ProductTypeID != nil {
		product.ProductTypeID = *dto.ProductTypeID
	}
	if dto.Name != nil {
		product.Name = *dto.Name
	}
	if dto.Description != nil {
		product.Description = dto.Description
	}
	if dto.MinPartnerPrice != nil {
		product.MinPartnerPrice = *dto.MinPartnerPrice
	}
	if dto.RollWidth != nil {
		product.RollWidth = dto.RollWidth
	}

	return product
}
