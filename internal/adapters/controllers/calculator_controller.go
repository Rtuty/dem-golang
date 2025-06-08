package controllers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/adapters/controllers/dto"
	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/usecases"

	"github.com/gin-gonic/gin"
)

// CalculatorController обрабатывает HTTP запросы для калькулятора
type CalculatorController struct {
	calculatorUseCase *usecases.CalculatorUseCase
	materialUseCase   *usecases.MaterialUseCase
	productUseCase    *usecases.ProductUseCase
}

// NewCalculatorController создает новый контроллер калькулятора
func NewCalculatorController(
	calculatorUseCase *usecases.CalculatorUseCase,
	materialUseCase *usecases.MaterialUseCase,
	productUseCase *usecases.ProductUseCase,
) *CalculatorController {
	return &CalculatorController{
		calculatorUseCase: calculatorUseCase,
		materialUseCase:   materialUseCase,
		productUseCase:    productUseCase,
	}
}

// GetCalculatorPage отображает страницу калькулятора
func (c *CalculatorController) GetCalculatorPage(ctx *gin.Context) {
	// Получаем типы продукции и материалов для выпадающих списков
	productTypes, err := c.productUseCase.GetProductTypes()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения типов продукции",
		})
		return
	}

	materialTypes, err := c.materialUseCase.GetMaterialTypes()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения типов материалов",
		})
		return
	}

	ctx.HTML(http.StatusOK, "calculator.html", gin.H{
		"title":         "Калькулятор материалов",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
	})
}

// CalculateMaterial обрабатывает форму расчета материала
func (c *CalculatorController) CalculateMaterial(ctx *gin.Context) {
	// Получаем данные из формы
	productTypeID, err := strconv.Atoi(ctx.PostForm("product_type_id"))
	if err != nil {
		c.renderCalculatorWithError(ctx, "Некорректный тип продукции")
		return
	}

	materialTypeID, err := strconv.Atoi(ctx.PostForm("material_type_id"))
	if err != nil {
		c.renderCalculatorWithError(ctx, "Некорректный тип материала")
		return
	}

	productQuantity, err := strconv.Atoi(ctx.PostForm("product_quantity"))
	if err != nil {
		c.renderCalculatorWithError(ctx, "Некорректное количество продукции")
		return
	}

	productParam1, err := strconv.ParseFloat(ctx.PostForm("product_param1"), 64)
	if err != nil {
		c.renderCalculatorWithError(ctx, "Некорректный первый параметр")
		return
	}

	productParam2, err := strconv.ParseFloat(ctx.PostForm("product_param2"), 64)
	if err != nil {
		c.renderCalculatorWithError(ctx, "Некорректный второй параметр")
		return
	}

	materialInStock, err := strconv.ParseFloat(ctx.PostForm("material_in_stock"), 64)
	if err != nil {
		materialInStock = 0
	}

	// Создаем запрос на расчет
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   productTypeID,
		MaterialTypeID:  materialTypeID,
		ProductQuantity: productQuantity,
		ProductParam1:   productParam1,
		ProductParam2:   productParam2,
		MaterialInStock: materialInStock,
	}

	// Выполняем расчет
	result, err := c.calculatorUseCase.CalculateRequiredMaterial(request)
	if err != nil {
		c.renderCalculatorWithError(ctx, err.Error())
		return
	}

	// Получаем данные для отображения
	productTypes, _ := c.productUseCase.GetProductTypes()
	materialTypes, _ := c.materialUseCase.GetMaterialTypes()

	ctx.HTML(http.StatusOK, "calculator.html", gin.H{
		"title":         "Калькулятор материалов",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
		"result":        result,
		"request":       request,
	})
}

// CalculateMaterialAPI обрабатывает API запрос на расчет материала
func (c *CalculatorController) CalculateMaterialAPI(ctx *gin.Context) {
	var requestDTO dto.MaterialCalculationRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректные данные запроса",
		})
		return
	}

	// Преобразуем в доменную сущность
	request := requestDTO.ToEntity()

	// Выполняем расчет
	result, err := c.calculatorUseCase.CalculateRequiredMaterial(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"required_quantity": result,
	})
}

// renderCalculatorWithError отображает страницу калькулятора с ошибкой
func (c *CalculatorController) renderCalculatorWithError(ctx *gin.Context, errorMsg string) {
	productTypes, _ := c.productUseCase.GetProductTypes()
	materialTypes, _ := c.materialUseCase.GetMaterialTypes()

	ctx.HTML(http.StatusBadRequest, "calculator.html", gin.H{
		"title":         "Калькулятор материалов",
		"productTypes":  productTypes,
		"materialTypes": materialTypes,
		"error":         errorMsg,
	})
}
