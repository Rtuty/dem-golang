package handlers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MaterialHandler обрабатывает HTTP запросы для материалов
type MaterialHandler struct {
	materialRepo repository.MaterialRepository
	logger       *zap.SugaredLogger
}

// NewMaterialHandler создает новый обработчик материалов
func NewMaterialHandler(materialRepo repository.MaterialRepository, logger *zap.SugaredLogger) *MaterialHandler {
	return &MaterialHandler{
		materialRepo: materialRepo,
		logger:       logger,
	}
}

// CreateMaterialRequest представляет запрос на создание материала
type CreateMaterialRequest struct {
	Article              string  `json:"article" binding:"required"`
	MaterialTypeID       int     `json:"material_type_id" binding:"required"`
	Name                 string  `json:"name" binding:"required"`
	Description          *string `json:"description"`
	MeasurementUnitID    int     `json:"measurement_unit_id" binding:"required"`
	PackageQuantity      float64 `json:"package_quantity" binding:"required,gt=0"`
	CostPerUnit          float64 `json:"cost_per_unit" binding:"required,gt=0"`
	StockQuantity        float64 `json:"stock_quantity" binding:"required,gte=0"`
	MinStockQuantity     float64 `json:"min_stock_quantity" binding:"required,gte=0"`
	ImagePath            *string `json:"image_path"`
}

// MaterialResponse представляет ответ с информацией о материале
type MaterialResponse struct {
	ID                   int     `json:"id"`
	Article              string  `json:"article"`
	MaterialTypeID       int     `json:"material_type_id"`
	Name                 string  `json:"name"`
	Description          *string `json:"description"`
	MeasurementUnitID    int     `json:"measurement_unit_id"`
	PackageQuantity      float64 `json:"package_quantity"`
	CostPerUnit          float64 `json:"cost_per_unit"`
	StockQuantity        float64 `json:"stock_quantity"`
	MinStockQuantity     float64 `json:"min_stock_quantity"`
	ImagePath            *string `json:"image_path"`
}

// CreateMaterial создает новый материал
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	var req CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Ошибка валидации запроса",
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	// Создаем entity материала
	material, err := entity.NewMaterial(
		req.Article,
		entity.ID(req.MaterialTypeID),
		req.Name,
		entity.ID(req.MeasurementUnitID),
		req.PackageQuantity,
		entity.Money(req.CostPerUnit),
	)
	if err != nil {
		h.logger.Errorw("Ошибка создания entity материала",
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные материала"})
		return
	}

	// Создаем материал через репозиторий
	err = h.materialRepo.Create(c.Request.Context(), material)
	if err != nil {
		h.logger.Errorw("Ошибка создания материала",
			"error", err,
			"article", req.Article,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания материала"})
		return
	}

	h.logger.Infow("Материал успешно создан",
		"article", req.Article,
		"request_id", c.GetString("request_id"),
	)

	c.JSON(http.StatusCreated, gin.H{"message": "Материал успешно создан"})
}

// GetMaterial возвращает материал по ID
func (h *MaterialHandler) GetMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Неверный ID материала",
			"id", idStr,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID материала"})
		return
	}

	material, err := h.materialRepo.GetByID(c.Request.Context(), entity.ID(id))
	if err != nil {
		h.logger.Errorw("Ошибка получения материала",
			"id", id,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "Материал не найден"})
		return
	}

	response := h.materialToResponse(material)
	c.JSON(http.StatusOK, response)
}

// GetMaterials возвращает список материалов
func (h *MaterialHandler) GetMaterials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	materials, err := h.materialRepo.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Errorw("Ошибка получения списка материалов",
			"page", page,
			"limit", limit,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка материалов"})
		return
	}

	response := make([]MaterialResponse, len(materials))
	for i, material := range materials {
		response[i] = h.materialToResponse(material)
	}

	c.JSON(http.StatusOK, gin.H{
		"materials": response,
		"page":      page,
		"limit":     limit,
	})
}

// SearchMaterials ищет материалы по запросу
func (h *MaterialHandler) SearchMaterials(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поисковый запрос не может быть пустым"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	materials, err := h.materialRepo.Search(c.Request.Context(), query, limit, offset)
	if err != nil {
		h.logger.Errorw("Ошибка поиска материалов",
			"query", query,
			"page", page,
			"limit", limit,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска материалов"})
		return
	}

	response := make([]MaterialResponse, len(materials))
	for i, material := range materials {
		response[i] = h.materialToResponse(material)
	}

	c.JSON(http.StatusOK, gin.H{
		"materials": response,
		"query":     query,
		"page":      page,
		"limit":     limit,
	})
}

// GetLowStockMaterials возвращает материалы с низким остатком
func (h *MaterialHandler) GetLowStockMaterials(c *gin.Context) {
	materials, err := h.materialRepo.GetLowStockMaterials(c.Request.Context())
	if err != nil {
		h.logger.Errorw("Ошибка получения материалов с низким остатком",
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения материалов с низким остатком"})
		return
	}

	response := make([]MaterialResponse, len(materials))
	for i, material := range materials {
		response[i] = h.materialToResponse(material)
	}

	c.JSON(http.StatusOK, gin.H{"materials": response})
}

// materialToResponse преобразует entity.Material в MaterialResponse
func (h *MaterialHandler) materialToResponse(material *entity.Material) MaterialResponse {
	return MaterialResponse{
		ID:                int(material.ID()),
		Article:           material.Article(),
		MaterialTypeID:    int(material.MaterialTypeID()),
		Name:              material.Name(),
		MeasurementUnitID: int(material.MeasurementUnitID()),
		PackageQuantity:   material.PackageQuantity(),
		CostPerUnit:       float64(material.CostPerUnit()),
		StockQuantity:     material.StockQuantity(),
		MinStockQuantity:  material.MinStockQuantity(),
	}
} 