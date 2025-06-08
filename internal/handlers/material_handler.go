package handlers

import (
	"net/http"

	"wallpaper-system/internal/models"
	"wallpaper-system/internal/services"

	"github.com/gin-gonic/gin"
)

// MaterialHandler представляет хендлер для работы с материалами
type MaterialHandler struct {
	materialService *services.MaterialService
	productService  *services.ProductService
}

// NewMaterialHandler создает новый хендлер материалов
func NewMaterialHandler(materialService *services.MaterialService, productService *services.ProductService) *MaterialHandler {
	return &MaterialHandler{
		materialService: materialService,
		productService:  productService,
	}
}

// GetMaterials возвращает список всех материалов
func (h *MaterialHandler) GetMaterials(c *gin.Context) {
	materials, err := h.materialService.GetAllMaterials()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки списка материалов: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "materials.html", gin.H{
		"title":     "Список материалов - Наш декор",
		"materials": materials,
	})
}

// ShowCalculatorForm показывает форму калькулятора материалов
func (h *MaterialHandler) ShowCalculatorForm(c *gin.Context) {
	// Получаем типы продукции
	productTypes, err := h.productService.GetProductTypes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки типов продукции: " + err.Error(),
		})
		return
	}

	// Получаем типы материалов
	materialTypes, err := h.materialService.GetMaterialTypes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Ошибка",
			"error": "Ошибка загрузки типов материалов: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "calculator.html", gin.H{
		"title":         "Калькулятор материалов - Наш декор",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
	})
}

// CalculateMaterial выполняет расчет необходимого количества материала
func (h *MaterialHandler) CalculateMaterial(c *gin.Context) {
	var req models.MaterialCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка валидации данных: " + err.Error(),
		})
		return
	}

	result, err := h.materialService.CalculateRequiredMaterial(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка расчета: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CalculateMaterialForm обрабатывает форму расчета материалов
func (h *MaterialHandler) CalculateMaterialForm(c *gin.Context) {
	var req models.MaterialCalculationRequest
	if err := c.ShouldBind(&req); err != nil {
		// Получаем данные для повторного отображения формы
		productTypes, _ := h.productService.GetProductTypes()
		materialTypes, _ := h.materialService.GetMaterialTypes()

		c.HTML(http.StatusBadRequest, "calculator.html", gin.H{
			"title":         "Калькулятор материалов - Наш декор",
			"productTypes":  productTypes,
			"materialTypes": materialTypes,
			"error":         "Ошибка валидации данных: " + err.Error(),
			"request":       req,
		})
		return
	}

	result, err := h.materialService.CalculateRequiredMaterial(&req)
	if err != nil {
		// Получаем данные для повторного отображения формы
		productTypes, _ := h.productService.GetProductTypes()
		materialTypes, _ := h.materialService.GetMaterialTypes()

		c.HTML(http.StatusInternalServerError, "calculator.html", gin.H{
			"title":         "Калькулятор материалов - Наш декор",
			"productTypes":  productTypes,
			"materialTypes": materialTypes,
			"error":         "Ошибка расчета: " + err.Error(),
			"request":       req,
		})
		return
	}

	// Получаем данные для отображения результата
	productTypes, _ := h.productService.GetProductTypes()
	materialTypes, _ := h.materialService.GetMaterialTypes()

	c.HTML(http.StatusOK, "calculator.html", gin.H{
		"title":         "Калькулятор материалов - Наш декор",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
		"request":       req,
		"result":        result,
		"success":       true,
	})
}
