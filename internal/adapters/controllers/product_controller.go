package controllers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/adapters/controllers/dto"
	"wallpaper-system/internal/usecases"

	"github.com/gin-gonic/gin"
)

// ProductController обрабатывает HTTP запросы для продукции
type ProductController struct {
	productUseCase  *usecases.ProductUseCase
	materialUseCase *usecases.MaterialUseCase
}

// NewProductController создает новый контроллер продукции
func NewProductController(
	productUseCase *usecases.ProductUseCase,
	materialUseCase *usecases.MaterialUseCase,
) *ProductController {
	return &ProductController{
		productUseCase:  productUseCase,
		materialUseCase: materialUseCase,
	}
}

// GetProductsPage отображает страницу со списком продукции
func (c *ProductController) GetProductsPage(ctx *gin.Context) {
	products, err := c.productUseCase.GetAllProducts()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения списка продукции",
		})
		return
	}

	// Преобразуем в DTO для отображения
	productDTOs := make([]dto.ProductListItemDTO, len(products))
	for i, product := range products {
		productDTOs[i] = dto.FromProductEntity(&product)
	}

	ctx.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Список продукции",
		"products": productDTOs,
	})
}

// GetProductDetailsPage отображает страницу с деталями продукции
func (c *ProductController) GetProductDetailsPage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID продукции",
		})
		return
	}

	product, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Продукция не найдена",
		})
		return
	}

	// Получаем материалы для продукции
	materials, err := c.materialUseCase.GetMaterialsForProduct(id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения материалов",
		})
		return
	}

	productDTO := dto.FromProductEntityWithMaterials(product, materials)

	ctx.HTML(http.StatusOK, "product_details.html", gin.H{
		"title":     "Детали продукции",
		"product":   productDTO,
		"materials": materials,
	})
}

// GetProducts возвращает список продукции в JSON
func (c *ProductController) GetProducts(ctx *gin.Context) {
	products, err := c.productUseCase.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения списка продукции",
		})
		return
	}

	// Преобразуем в DTO
	productDTOs := make([]dto.ProductListItemDTO, len(products))
	for i, product := range products {
		productDTOs[i] = dto.FromProductEntity(&product)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": productDTOs,
	})
}

// GetProductByID возвращает продукцию по ID в JSON
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный ID продукции",
		})
		return
	}

	product, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Продукция не найдена",
		})
		return
	}

	productDTO := dto.FromProductEntity(product)
	ctx.JSON(http.StatusOK, productDTO)
}

// CreateProduct создает новую продукцию
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var request dto.CreateProductRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректные данные запроса",
		})
		return
	}

	// Преобразуем DTO в доменную сущность
	product := request.ToEntity()

	err := c.productUseCase.CreateProduct(product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Продукция успешно создана",
		"id":      product.ID,
	})
}

// UpdateProduct обновляет продукцию
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный ID продукции",
		})
		return
	}

	var request dto.UpdateProductRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректные данные запроса",
		})
		return
	}

	// Получаем существующую продукцию
	product, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Продукция не найдена",
		})
		return
	}

	// Обновляем поля
	request.UpdateEntity(product)

	err = c.productUseCase.UpdateProduct(product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Продукция успешно обновлена",
	})
}

// DeleteProduct удаляет продукцию
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный ID продукции",
		})
		return
	}

	err = c.productUseCase.DeleteProduct(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Продукция успешно удалена",
	})
}
