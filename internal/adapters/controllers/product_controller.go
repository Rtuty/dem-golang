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
	productUseCase  usecases.ProductUseCaseInterface
	materialUseCase usecases.MaterialUseCaseInterface
}

// NewProductController создает новый контроллер продукции
func NewProductController(
	productUseCase usecases.ProductUseCaseInterface,
	materialUseCase usecases.MaterialUseCaseInterface,
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

	ctx.HTML(http.StatusOK, "product_detail.html", gin.H{
		"title":     "Детали продукции",
		"product":   productDTO,
		"materials": materials,
	})
}

// GetProducts возвращает список продукции в JSON
func (c *ProductController) GetProducts(ctx *gin.Context) {
	products, err := c.productUseCase.GetAllProducts()
	if err != nil {
		response := dto.NewErrorResponse("Ошибка получения списка продукции")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// Преобразуем в DTO
	productDTOs := make([]dto.ProductListItemDTO, len(products))
	for i, product := range products {
		productDTOs[i] = dto.FromProductEntity(&product)
	}

	response := dto.NewSuccessResponse("Список продукции получен", productDTOs)
	ctx.JSON(http.StatusOK, response)
}

// GetProductByID возвращает продукцию по ID в JSON
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := dto.NewErrorResponse("Некорректный ID продукции")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		response := dto.NewErrorResponse("Продукция не найдена")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	productDTO := dto.FromProductEntity(product)
	response := dto.NewSuccessResponse("Продукция найдена", productDTO)
	ctx.JSON(http.StatusOK, response)
}

// CreateProduct создает новую продукцию
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var request dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response := dto.NewErrorResponse("Некорректные данные запроса")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// Преобразуем DTO в доменную сущность
	product := request.ToEntity()

	err := c.productUseCase.CreateProduct(product)
	if err != nil {
		response := dto.NewErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse("Продукция успешно создана", gin.H{"id": product.ID})
	ctx.JSON(http.StatusCreated, response)
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

	var request dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректные данные запроса",
		})
		return
	}

	// Преобразуем DTO в доменную сущность
	product := request.ToEntity(id)

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
		response := dto.NewErrorResponse("Некорректный ID продукции")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.productUseCase.DeleteProduct(id)
	if err != nil {
		response := dto.NewErrorResponse(err.Error())
		ctx.JSON(http.StatusNotFound, response) // Изменил на 404 если продукция не найдена
		return
	}

	response := dto.NewSuccessResponse("Продукция успешно удалена", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetCreateProductPage отображает форму создания новой продукции
func (c *ProductController) GetCreateProductPage(ctx *gin.Context) {
	// Получаем типы продукции для выпадающего списка
	productTypes, err := c.productUseCase.GetProductTypes()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки типов продукции: " + err.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":        "Создание продукции",
		"isEdit":       false,
		"formAction":   "/products",
		"product":      nil,
		"productTypes": productTypes,
	})
}

// CreateProductWeb обрабатывает создание продукции через веб-форму
func (c *ProductController) CreateProductWeb(ctx *gin.Context) {
	var request dto.CreateProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		// Получаем типы продукции для формы в случае ошибки
		productTypes, _ := c.productUseCase.GetProductTypes()
		
		ctx.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Создание продукции",
			"isEdit":       false,
			"formAction":   "/products",
			"error":        "Некорректные данные формы: " + err.Error(),
			"product":      request,
			"productTypes": productTypes,
		})
		return
	}

	// Преобразуем DTO в доменную сущность
	product := request.ToEntity()

	err := c.productUseCase.CreateProduct(product)
	if err != nil {
		// Получаем типы продукции для формы в случае ошибки
		productTypes, _ := c.productUseCase.GetProductTypes()
		
		ctx.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Создание продукции",
			"isEdit":       false,
			"formAction":   "/products",
			"error":        "Ошибка создания продукции: " + err.Error(),
			"product":      request,
			"productTypes": productTypes,
		})
		return
	}

	// Успешное создание - редирект на список продукции
	ctx.Redirect(http.StatusFound, "/products")
}

// GetEditProductPage отображает форму редактирования продукции
func (c *ProductController) GetEditProductPage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID продукции",
		})
		return
	}

	// Получаем продукцию для редактирования
	product, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Продукция не найдена: " + err.Error(),
		})
		return
	}

	// Получаем типы продукции для выпадающего списка
	productTypes, err := c.productUseCase.GetProductTypes()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки типов продукции: " + err.Error(),
		})
		return
	}

	// Преобразуем в DTO для отображения в форме
	productDTO := dto.FromProductEntity(product)

	ctx.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":        "Редактирование продукции",
		"isEdit":       true,
		"formAction":   "/products/" + strconv.Itoa(id),
		"product":      productDTO,
		"productTypes": productTypes,
	})
}

// UpdateProductWeb обрабатывает обновление продукции через веб-форму
func (c *ProductController) UpdateProductWeb(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID продукции",
		})
		return
	}

	var request dto.UpdateProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		// Получаем продукцию и типы продукции для формы в случае ошибки
		product, _ := c.productUseCase.GetProductByID(id)
		productTypes, _ := c.productUseCase.GetProductTypes()
		productDTO := dto.FromProductEntity(product)
		
		ctx.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Редактирование продукции",
			"isEdit":       true,
			"formAction":   "/products/" + strconv.Itoa(id),
			"error":        "Некорректные данные формы: " + err.Error(),
			"product":      productDTO,
			"productTypes": productTypes,
		})
		return
	}

	// Преобразуем DTO в доменную сущность
	product := request.ToEntity(id)

	err = c.productUseCase.UpdateProduct(product)
	if err != nil {
		// Получаем продукцию и типы продукции для формы в случае ошибки
		existingProduct, _ := c.productUseCase.GetProductByID(id)
		productTypes, _ := c.productUseCase.GetProductTypes()
		productDTO := dto.FromProductEntity(existingProduct)
		
		ctx.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Редактирование продукции",
			"isEdit":       true,
			"formAction":   "/products/" + strconv.Itoa(id),
			"error":        "Ошибка обновления продукции: " + err.Error(),
			"product":      productDTO,
			"productTypes": productTypes,
		})
		return
	}

	// Успешное обновление - редирект на детали продукции
	ctx.Redirect(http.StatusFound, "/products/"+strconv.Itoa(id))
}
