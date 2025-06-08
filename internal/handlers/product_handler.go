package handlers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/models"
	"wallpaper-system/internal/services"

	"github.com/gin-gonic/gin"
)

// ProductHandler представляет хендлер для работы с продукцией
type ProductHandler struct {
	productService  *services.ProductService
	materialService *services.MaterialService
}

// NewProductHandler создает новый хендлер продукции
func NewProductHandler(productService *services.ProductService, materialService *services.MaterialService) *ProductHandler {
	return &ProductHandler{
		productService:  productService,
		materialService: materialService,
	}
}

// GetProducts возвращает список продукции (главная страница)
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки списка продукции: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Список продукции - Наш декор",
		"products": products,
	})
}

// GetProduct возвращает детальную информацию о продукции
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Неверный ID продукции",
		})
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Продукция не найдена: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "product_detail.html", gin.H{
		"title":   "Детали продукции - " + product.Name,
		"product": product,
	})
}

// GetProductMaterials возвращает список материалов для продукции
func (h *ProductHandler) GetProductMaterials(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Неверный ID продукции",
		})
		return
	}

	// Получаем информацию о продукции
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Продукция не найдена: " + err.Error(),
		})
		return
	}

	// Получаем материалы для продукции
	materials, err := h.materialService.GetMaterialsForProduct(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки материалов: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "product_materials.html", gin.H{
		"title":     "Материалы для " + product.Name,
		"product":   product,
		"materials": materials,
	})
}

// ShowCreateProductForm показывает форму создания продукции
func (h *ProductHandler) ShowCreateProductForm(c *gin.Context) {
	productTypes, err := h.productService.GetProductTypes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки типов продукции: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":        "Добавить продукцию - Наш декор",
		"product":      nil,
		"productTypes": productTypes,
		"isEdit":       false,
	})
}

// ShowEditProductForm показывает форму редактирования продукции
func (h *ProductHandler) ShowEditProductForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Неверный ID продукции",
		})
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Продукция не найдена: " + err.Error(),
		})
		return
	}

	productTypes, err := h.productService.GetProductTypes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки типов продукции: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":        "Редактировать продукцию - " + product.Name,
		"product":      product,
		"productTypes": productTypes,
		"isEdit":       true,
	})
}

// CreateProduct создает новую продукцию
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		productTypes, _ := h.productService.GetProductTypes()
		c.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Добавить продукцию - Наш декор",
			"product":      nil,
			"productTypes": productTypes,
			"isEdit":       false,
			"error":        "Ошибка валидации данных: " + err.Error(),
		})
		return
	}

	_, err := h.productService.CreateProduct(&req)
	if err != nil {
		productTypes, _ := h.productService.GetProductTypes()
		c.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Добавить продукцию - Наш декор",
			"product":      nil,
			"productTypes": productTypes,
			"isEdit":       false,
			"error":        "Ошибка создания продукции: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

// UpdateProduct обновляет продукцию
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Неверный ID продукции",
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		product, _ := h.productService.GetProductByID(id)
		productTypes, _ := h.productService.GetProductTypes()
		c.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Редактировать продукцию",
			"product":      product,
			"productTypes": productTypes,
			"isEdit":       true,
			"error":        "Ошибка валидации данных: " + err.Error(),
		})
		return
	}

	err = h.productService.UpdateProduct(id, &req)
	if err != nil {
		product, _ := h.productService.GetProductByID(id)
		productTypes, _ := h.productService.GetProductTypes()
		c.HTML(http.StatusBadRequest, "product_form.html", gin.H{
			"title":        "Редактировать продукцию",
			"product":      product,
			"productTypes": productTypes,
			"isEdit":       true,
			"error":        "Ошибка обновления продукции: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

// DeleteProduct удаляет продукцию
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукции"})
		return
	}

	err = h.productService.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Продукция успешно удалена"})
}
