package usecase

import (
	"context"
	"fmt"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/service"
)

// ProductUseCase представляет вариант использования для продукции
type ProductUseCase struct {
	productService *service.ProductService
}

// NewProductUseCase создает новый use case для продукции
func NewProductUseCase(productService *service.ProductService) *ProductUseCase {
	return &ProductUseCase{
		productService: productService,
	}
}

// CreateProductRequest представляет запрос на создание продукции
type CreateProductRequest struct {
	Article         string  `json:"article" validate:"required,min=1,max=100"`
	ProductTypeID   int     `json:"product_type_id" validate:"required,min=1"`
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	Description     *string `json:"description" validate:"omitempty,max=1000"`
	MinPartnerPrice float64 `json:"min_partner_price" validate:"required,min=0"`
	RollWidth       *float64 `json:"roll_width" validate:"omitempty,min=0"`
}

// UpdateProductRequest представляет запрос на обновление продукции
type UpdateProductRequest struct {
	Name            *string  `json:"name" validate:"omitempty,min=1,max=255"`
	Description     *string  `json:"description" validate:"omitempty,max=1000"`
	MinPartnerPrice *float64 `json:"min_partner_price" validate:"omitempty,min=0"`
	RollWidth       *float64 `json:"roll_width" validate:"omitempty,min=0"`
}

// ProductResponse представляет ответ с данными продукции
type ProductResponse struct {
	ID              int                        `json:"id"`
	Article         string                     `json:"article"`
	ProductTypeID   int                        `json:"product_type_id"`
	Name            string                     `json:"name"`
	Description     *string                    `json:"description"`
	ImagePath       *string                    `json:"image_path"`
	MinPartnerPrice float64                    `json:"min_partner_price"`
	RollWidth       *float64                   `json:"roll_width"`
	CreatedAt       string                     `json:"created_at"`
	UpdatedAt       string                     `json:"updated_at"`
	ProductType     *ProductTypeResponse       `json:"product_type,omitempty"`
	Materials       []ProductMaterialResponse  `json:"materials,omitempty"`
	CalculatedPrice *float64                   `json:"calculated_price,omitempty"`
}

// ProductTypeResponse представляет ответ с данными типа продукции
type ProductTypeResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Coefficient float64 `json:"coefficient"`
}

// ProductMaterialResponse представляет ответ с данными материала продукции
type ProductMaterialResponse struct {
	ID              int              `json:"id"`
	MaterialID      int              `json:"material_id"`
	QuantityPerUnit float64          `json:"quantity_per_unit"`
	Material        *MaterialResponse `json:"material,omitempty"`
}

// MaterialResponse представляет ответ с данными материала
type MaterialResponse struct {
	ID          int     `json:"id"`
	Article     string  `json:"article"`
	Name        string  `json:"name"`
	CostPerUnit float64 `json:"cost_per_unit"`
}

// ProductListResponse представляет ответ со списком продукции
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// CreateProduct создает новую продукцию
func (uc *ProductUseCase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	product, err := uc.productService.CreateProduct(
		ctx,
		req.Article,
		entity.ID(req.ProductTypeID),
		req.Name,
		entity.Money(req.MinPartnerPrice),
		req.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания продукции: %w", err)
	}

	return uc.entityToResponse(product), nil
}

// GetProductByID возвращает продукцию по ID
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id int) (*ProductResponse, error) {
	product, err := uc.productService.GetProductByID(ctx, entity.ID(id))
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	return uc.entityToResponse(product), nil
}

// GetAllProducts возвращает список всех продукций с пагинацией
func (uc *ProductUseCase) GetAllProducts(ctx context.Context, limit, offset int) (*ProductListResponse, error) {
	// Устанавливаем ограничения по умолчанию
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	products, total, err := uc.productService.GetAllProducts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка продукции: %w", err)
	}

	response := &ProductListResponse{
		Products: make([]ProductResponse, 0, len(products)),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	for _, product := range products {
		response.Products = append(response.Products, *uc.entityToResponse(product))
	}

	return response, nil
}

// UpdateProduct обновляет продукцию
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id int, req *UpdateProductRequest) error {
	var minPartnerPrice *entity.Money
	if req.MinPartnerPrice != nil {
		price := entity.Money(*req.MinPartnerPrice)
		minPartnerPrice = &price
	}

	err := uc.productService.UpdateProduct(
		ctx,
		entity.ID(id),
		req.Name,
		req.Description,
		minPartnerPrice,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления продукции: %w", err)
	}

	return nil
}

// DeleteProduct удаляет продукцию
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id int) error {
	err := uc.productService.DeleteProduct(ctx, entity.ID(id))
	if err != nil {
		return fmt.Errorf("ошибка удаления продукции: %w", err)
	}

	return nil
}

// AddMaterialToProduct добавляет материал к продукции
func (uc *ProductUseCase) AddMaterialToProduct(ctx context.Context, productID, materialID int, quantityPerUnit float64) error {
	err := uc.productService.AddMaterialToProduct(
		ctx,
		entity.ID(productID),
		entity.ID(materialID),
		quantityPerUnit,
	)
	if err != nil {
		return fmt.Errorf("ошибка добавления материала к продукции: %w", err)
	}

	return nil
}

// RemoveMaterialFromProduct удаляет материал из продукции
func (uc *ProductUseCase) RemoveMaterialFromProduct(ctx context.Context, productID, materialID int) error {
	err := uc.productService.RemoveMaterialFromProduct(
		ctx,
		entity.ID(productID),
		entity.ID(materialID),
	)
	if err != nil {
		return fmt.Errorf("ошибка удаления материала из продукции: %w", err)
	}

	return nil
}

// CalculateProductPrice рассчитывает цену продукции
func (uc *ProductUseCase) CalculateProductPrice(ctx context.Context, productID int) (float64, error) {
	price, err := uc.productService.CalculateProductPrice(ctx, entity.ID(productID))
	if err != nil {
		return 0, fmt.Errorf("ошибка расчета цены продукции: %w", err)
	}

	return float64(price), nil
}

// entityToResponse преобразует доменную сущность в ответ
func (uc *ProductUseCase) entityToResponse(product *entity.Product) *ProductResponse {
	response := &ProductResponse{
		ID:              int(product.ID()),
		Article:         product.Article(),
		ProductTypeID:   int(product.ProductTypeID()),
		Name:            product.Name(),
		Description:     product.Description(),
		ImagePath:       product.ImagePath(),
		MinPartnerPrice: float64(product.MinPartnerPrice()),
		RollWidth:       product.RollWidth(),
		CreatedAt:       product.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       product.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	if product.CalculatedPrice() != nil {
		price := float64(*product.CalculatedPrice())
		response.CalculatedPrice = &price
	}

	if product.ProductType() != nil {
		response.ProductType = &ProductTypeResponse{
			ID:          int(product.ProductType().ID()),
			Name:        product.ProductType().Name(),
			Coefficient: product.ProductType().Coefficient(),
		}
	}

	if len(product.Materials()) > 0 {
		response.Materials = make([]ProductMaterialResponse, 0, len(product.Materials()))
		for _, pm := range product.Materials() {
			materialResp := ProductMaterialResponse{
				ID:              int(pm.ID()),
				MaterialID:      int(pm.MaterialID()),
				QuantityPerUnit: pm.QuantityPerUnit(),
			}

			if pm.Material() != nil {
				materialResp.Material = &MaterialResponse{
					ID:          int(pm.Material().ID()),
					Article:     pm.Material().Article(),
					Name:        pm.Material().Name(),
					CostPerUnit: float64(pm.Material().CostPerUnit()),
				}
			}

			response.Materials = append(response.Materials, materialResp)
		}
	}

	return response
} 