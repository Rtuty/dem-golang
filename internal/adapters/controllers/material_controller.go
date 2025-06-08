package controllers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/usecases"

	"github.com/gin-gonic/gin"
)

// MaterialController обрабатывает HTTP запросы для материалов
type MaterialController struct {
	materialUseCase *usecases.MaterialUseCase
}

// NewMaterialController создает новый контроллер материалов
func NewMaterialController(materialUseCase *usecases.MaterialUseCase) *MaterialController {
	return &MaterialController{
		materialUseCase: materialUseCase,
	}
}

// GetMaterialsPage отображает страницу со списком материалов
func (mc *MaterialController) GetMaterialsPage(c *gin.Context) {
	materials, err := mc.materialUseCase.GetAllMaterials()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки материалов: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "materials.html", gin.H{
		"title":     "Материалы",
		"materials": materials,
	})
}

// GetMaterialsAPI возвращает список материалов в формате JSON
func (mc *MaterialController) GetMaterials(c *gin.Context) {
	materials, err := mc.materialUseCase.GetAllMaterials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка загрузки материалов: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"materials": materials,
	})
}

// GetMaterialByID возвращает материал по ID в формате JSON
func (mc *MaterialController) GetMaterialByID(c *gin.Context) {
	id := c.Param("id")
	
	materialID, err := parseIDParam(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID материала"})
		return
	}

	material, err := mc.materialUseCase.GetMaterialByID(materialID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Материал не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"material": material,
	})
}

// parseIDParam парсит ID из строкового параметра
func parseIDParam(idParam string) (int, error) {
	return strconv.Atoi(idParam)
} 