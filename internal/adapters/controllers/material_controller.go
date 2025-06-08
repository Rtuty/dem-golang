package controllers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/adapters/controllers/dto"
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

// GetMaterialDetailsPage отображает страницу с деталями материала
func (mc *MaterialController) GetMaterialDetailsPage(c *gin.Context) {
	id := c.Param("id")
	
	materialID, err := parseIDParam(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID материала",
		})
		return
	}

	material, err := mc.materialUseCase.GetMaterialByID(materialID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Материал не найден: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "material_detail.html", gin.H{
		"title":    "Детали материала",
		"material": material,
	})
}

// GetCreateMaterialPage отображает страницу создания материала
func (mc *MaterialController) GetCreateMaterialPage(c *gin.Context) {
	// Получаем типы материалов
	materialTypes, err := mc.materialUseCase.GetMaterialTypes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки типов материалов: " + err.Error(),
		})
		return
	}

	// Получаем единицы измерения
	measurementUnits, err := mc.materialUseCase.GetMeasurementUnits()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки единиц измерения: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "material_form.html", gin.H{
		"title":            "Создание материала",
		"material":         nil,
		"materialTypes":    materialTypes,
		"measurementUnits": measurementUnits,
		"isEdit":           false,
	})
}

// CreateMaterialWeb создает новый материал через веб-форму
func (mc *MaterialController) CreateMaterialWeb(c *gin.Context) {
	var dto dto.CreateMaterialDTO
	
	if err := c.ShouldBind(&dto); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Ошибка обработки формы: " + err.Error(),
		})
		return
	}

	material := dto.ToEntity()
	err := mc.materialUseCase.CreateMaterial(material)
	if err != nil {
		// Получаем данные для повторного отображения формы
		materialTypes, _ := mc.materialUseCase.GetMaterialTypes()
		measurementUnits, _ := mc.materialUseCase.GetMeasurementUnits()
		
		c.HTML(http.StatusBadRequest, "material_form.html", gin.H{
			"title":            "Создание материала",
			"material":         material,
			"materialTypes":    materialTypes,
			"measurementUnits": measurementUnits,
			"isEdit":           false,
			"error":            err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/materials")
}

// GetEditMaterialPage отображает страницу редактирования материала
func (mc *MaterialController) GetEditMaterialPage(c *gin.Context) {
	id := c.Param("id")
	
	materialID, err := parseIDParam(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID материала",
		})
		return
	}

	material, err := mc.materialUseCase.GetMaterialByID(materialID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Материал не найден: " + err.Error(),
		})
		return
	}

	// Получаем типы материалов
	materialTypes, err := mc.materialUseCase.GetMaterialTypes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки типов материалов: " + err.Error(),
		})
		return
	}

	// Получаем единицы измерения
	measurementUnits, err := mc.materialUseCase.GetMeasurementUnits()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки единиц измерения: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "material_form.html", gin.H{
		"title":            "Редактирование материала",
		"material":         material,
		"materialTypes":    materialTypes,
		"measurementUnits": measurementUnits,
		"isEdit":           true,
	})
}

// UpdateMaterialWeb обновляет материал через веб-форму
func (mc *MaterialController) UpdateMaterialWeb(c *gin.Context) {
	id := c.Param("id")
	
	materialID, err := parseIDParam(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Некорректный ID материала",
		})
		return
	}

	var dto dto.UpdateMaterialDTO
	if err := c.ShouldBind(&dto); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Ошибка обработки формы: " + err.Error(),
		})
		return
	}

	material := dto.ToEntity()
	material.ID = materialID

	err = mc.materialUseCase.UpdateMaterial(material)
	if err != nil {
		// Получаем данные для повторного отображения формы
		materialTypes, _ := mc.materialUseCase.GetMaterialTypes()
		measurementUnits, _ := mc.materialUseCase.GetMeasurementUnits()
		
		c.HTML(http.StatusBadRequest, "material_form.html", gin.H{
			"title":            "Редактирование материала",
			"material":         material,
			"materialTypes":    materialTypes,
			"measurementUnits": measurementUnits,
			"isEdit":           true,
			"error":            err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/materials/"+strconv.Itoa(materialID))
}

// CreateMaterial создает новый материал через API
func (mc *MaterialController) CreateMaterial(c *gin.Context) {
	var dto dto.CreateMaterialDTO
	
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректные данные: " + err.Error(),
		})
		return
	}

	material := dto.ToEntity()
	err := mc.materialUseCase.CreateMaterial(material)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Материал успешно создан",
		"data":    gin.H{"id": material.ID},
	})
}

// UpdateMaterial обновляет материал через API
func (mc *MaterialController) UpdateMaterial(c *gin.Context) {
	id := c.Param("id")
	
	materialID, err := parseIDParam(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректный ID материала",
		})
		return
	}

	var dto dto.UpdateMaterialDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректные данные: " + err.Error(),
		})
		return
	}

	material := dto.ToEntity()
	material.ID = materialID

	err = mc.materialUseCase.UpdateMaterial(material)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Материал успешно обновлен",
	})
}

// DeleteMaterial удаляет материал через API
func (mc *MaterialController) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	
	materialID, err := parseIDParam(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректный ID материала",
		})
		return
	}

	err = mc.materialUseCase.DeleteMaterial(materialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Ошибка удаления материала: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Материал успешно удален",
	})
}

// GetMaterialTypes возвращает список типов материалов через API
func (mc *MaterialController) GetMaterialTypes(c *gin.Context) {
	materialTypes, err := mc.materialUseCase.GetMaterialTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения типов материалов: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": materialTypes,
	})
}

// GetMeasurementUnits возвращает список единиц измерения через API
func (mc *MaterialController) GetMeasurementUnits(c *gin.Context) {
	measurementUnits, err := mc.materialUseCase.GetMeasurementUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения единиц измерения: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": measurementUnits,
	})
}

// parseIDParam парсит ID из строкового параметра
func parseIDParam(idParam string) (int, error) {
	return strconv.Atoi(idParam)
} 