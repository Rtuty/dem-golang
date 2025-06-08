package handlers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/models"
	"wallpaper-system/internal/services"

	"github.com/gin-gonic/gin"
)

// CalculatorHandler представляет обработчик для расчетов
type CalculatorHandler struct {
	calculatorService *services.CalculatorService
	productService    *services.ProductService
	materialService   *services.MaterialService
}

// NewCalculatorHandler создает новый обработчик калькулятора
func NewCalculatorHandler(calculatorService *services.CalculatorService, productService *services.ProductService, materialService *services.MaterialService) *CalculatorHandler {
	return &CalculatorHandler{
		calculatorService: calculatorService,
		productService:    productService,
		materialService:   materialService,
	}
}

// ShowCalculatorForm показывает форму калькулятора материалов
func (h *CalculatorHandler) ShowCalculatorForm(c *gin.Context) {
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

// CalculateMaterialAmount выполняет расчет необходимого количества материала
func (h *CalculatorHandler) CalculateMaterialAmount(c *gin.Context) {
	// Получаем параметры из запроса
	productTypeID, err := strconv.Atoi(c.PostForm("product_type_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID типа продукции", "result": -1})
		return
	}

	materialTypeID, err := strconv.Atoi(c.PostForm("material_type_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID типа материала", "result": -1})
		return
	}

	productQuantity, err := strconv.Atoi(c.PostForm("product_quantity"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное количество продукции", "result": -1})
		return
	}

	productParam1, err := strconv.ParseFloat(c.PostForm("product_param1"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный параметр 1 продукции", "result": -1})
		return
	}

	productParam2, err := strconv.ParseFloat(c.PostForm("product_param2"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный параметр 2 продукции", "result": -1})
		return
	}

	materialInStock, err := strconv.ParseFloat(c.PostForm("material_in_stock"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное количество материала на складе", "result": -1})
		return
	}

	// Вызываем сервис для расчета
	result := h.calculatorService.CalculateRequiredMaterialAmount(
		productTypeID,
		materialTypeID,
		productQuantity,
		productParam1,
		productParam2,
		materialInStock,
	)

	if result == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Ошибка расчета. Проверьте корректность параметров.",
			"result": -1,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
		"message": "Расчет выполнен успешно",
	})
}

// CalculateMaterialForm обрабатывает форму расчета материалов
func (h *CalculatorHandler) CalculateMaterialForm(c *gin.Context) {
	// Получаем параметры из формы
	productTypeID, err := strconv.Atoi(c.PostForm("product_type_id"))
	if err != nil {
		renderCalculatorFormWithError(c, h, "Неверный ID типа продукции")
		return
	}

	materialTypeID, err := strconv.Atoi(c.PostForm("material_type_id"))
	if err != nil {
		renderCalculatorFormWithError(c, h, "Неверный ID типа материала")
		return
	}

	productQuantity, err := strconv.Atoi(c.PostForm("product_quantity"))
	if err != nil {
		renderCalculatorFormWithError(c, h, "Неверное количество продукции")
		return
	}

	productParam1, err := strconv.ParseFloat(c.PostForm("product_param1"), 64)
	if err != nil {
		renderCalculatorFormWithError(c, h, "Неверный параметр 1 продукции")
		return
	}

	productParam2, err := strconv.ParseFloat(c.PostForm("product_param2"), 64)
	if err != nil {
		renderCalculatorFormWithError(c, h, "Неверный параметр 2 продукции")
		return
	}

	materialInStock, err := strconv.ParseFloat(c.PostForm("material_in_stock"), 64)
	if err != nil {
		renderCalculatorFormWithError(c, h, "Неверное количество материала на складе")
		return
	}

	// Вызываем сервис для расчета
	result := h.calculatorService.CalculateRequiredMaterialAmount(
		productTypeID,
		materialTypeID,
		productQuantity,
		productParam1,
		productParam2,
		materialInStock,
	)

	if result == -1 {
		renderCalculatorFormWithError(c, h, "Ошибка расчета. Проверьте корректность параметров.")
		return
	}

	// Получаем данные для отображения результата
	productTypes, _ := h.productService.GetProductTypes()
	materialTypes, _ := h.materialService.GetMaterialTypes()

	// Формируем запрос для отображения в форме
	req := models.MaterialCalculationRequest{
		ProductTypeID:   productTypeID,
		MaterialTypeID:  materialTypeID,
		ProductQuantity: productQuantity,
		ProductParam1:   productParam1,
		ProductParam2:   productParam2,
		MaterialInStock: materialInStock,
	}

	c.HTML(http.StatusOK, "calculator.html", gin.H{
		"title":         "Калькулятор материалов - Наш декор",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
		"request":       req,
		"result":        result,
		"success":       true,
	})
}

// renderCalculatorFormWithError отображает форму калькулятора с сообщением об ошибке
func renderCalculatorFormWithError(c *gin.Context, h *CalculatorHandler, errorMessage string) {
	// Получаем данные для повторного отображения формы
	productTypes, _ := h.productService.GetProductTypes()
	materialTypes, _ := h.materialService.GetMaterialTypes()

	// Создаем запрос с данными из формы
	productTypeID, _ := strconv.Atoi(c.PostForm("product_type_id"))
	materialTypeID, _ := strconv.Atoi(c.PostForm("material_type_id"))
	productQuantity, _ := strconv.Atoi(c.PostForm("product_quantity"))
	productParam1, _ := strconv.ParseFloat(c.PostForm("product_param1"), 64)
	productParam2, _ := strconv.ParseFloat(c.PostForm("product_param2"), 64)
	materialInStock, _ := strconv.ParseFloat(c.PostForm("material_in_stock"), 64)

	req := models.MaterialCalculationRequest{
		ProductTypeID:   productTypeID,
		MaterialTypeID:  materialTypeID,
		ProductQuantity: productQuantity,
		ProductParam1:   productParam1,
		ProductParam2:   productParam2,
		MaterialInStock: materialInStock,
	}

	c.HTML(http.StatusBadRequest, "calculator.html", gin.H{
		"title":         "Калькулятор материалов - Наш декор",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
		"error":         errorMessage,
		"request":       req,
	})
}
