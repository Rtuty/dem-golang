package handlers

import (
	"net/http"
	"strconv"

	"wallpaper-system/internal/application/usecase"
	"wallpaper-system/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProductHandler обрабатывает HTTP запросы для продукции
type ProductHandler struct {
	productUsecase usecase.ProductUsecase
	logger         *zap.SugaredLogger
}

// NewProductHandler создает новый обработчик продукции
func NewProductHandler(productUsecase usecase.ProductUsecase, logger *zap.SugaredLogger) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		logger:         logger,
	}
}

// CreateProductRequest представляет запрос на создание продукции
type CreateProductRequest struct {
	Article           string  `json:"article" binding:"required"`
	ProductTypeID     int     `json:"product_type_id" binding:"required"`
	Name              string  `json:"name" binding:"required"`
	Description       *string `json:"description"`
	MinPartnerPrice   float64 `json:"min_partner_price" binding:"required,gt=0"`
	ImagePath         *string `json:"image_path"`
	WorkshopNumber    *string `json:"workshop_number"`
	RequiredWorkers   *int    `json:"required_workers"`
	RollWidth         *float64 `json:"roll_width"`
}

// UpdateProductRequest представляет запрос на обновление продукции
type UpdateProductRequest struct {
	Name              string   `json:"name" binding:"required"`
	Description       *string  `json:"description"`
	MinPartnerPrice   float64  `json:"min_partner_price" binding:"required,gt=0"`
	ImagePath         *string  `json:"image_path"`
	WorkshopNumber    *string  `json:"workshop_number"`
	RequiredWorkers   *int     `json:"required_workers"`
	RollWidth         *float64 `json:"roll_width"`
}

// ProductResponse представляет ответ с информацией о продукции
type ProductResponse struct {
	ID                int     `json:"id"`
	Article           string  `json:"article"`
	ProductTypeID     int     `json:"product_type_id"`
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	MinPartnerPrice   float64 `json:"min_partner_price"`
	ImagePath         *string `json:"image_path"`
	WorkshopNumber    *string `json:"workshop_number"`
	RequiredWorkers   *int    `json:"required_workers"`
	RollWidth         *float64 `json:"roll_width"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// CreateProduct создает новую продукцию
// @Summary Создать продукцию
// @Description Создает новую продукцию в системе
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Данные продукции"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Ошибка валидации запроса",
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	// Создаем entity продукции
	product, err := entity.NewProduct(
		req.Article,
		entity.ID(req.ProductTypeID),
		req.Name,
		entity.Money(req.MinPartnerPrice),
	)
	if err != nil {
		h.logger.Errorw("Ошибка создания entity продукции",
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные продукции"})
		return
	}

	// Создаем продукцию через usecase
	err = h.productUsecase.CreateProduct(c.Request.Context(), product)
	if err != nil {
		h.logger.Errorw("Ошибка создания продукции",
			"error", err,
			"article", req.Article,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания продукции"})
		return
	}

	h.logger.Infow("Продукция успешно создана",
		"article", req.Article,
		"request_id", c.GetString("request_id"),
	)

	c.JSON(http.StatusCreated, gin.H{"message": "Продукция успешно создана"})
}

// GetProduct возвращает продукцию по ID
// @Summary Получить продукцию
// @Description Возвращает информацию о продукции по ID
// @Tags products
// @Produce json
// @Param id path int true "ID продукции"
// @Success 200 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Неверный ID продукции",
			"id", idStr,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукции"})
		return
	}

	product, err := h.productUsecase.GetProductByID(c.Request.Context(), entity.ID(id))
	if err != nil {
		h.logger.Errorw("Ошибка получения продукции",
			"id", id,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукция не найдена"})
		return
	}

	response := h.productToResponse(product)
	c.JSON(http.StatusOK, response)
}

// GetProducts возвращает список продукций с пагинацией
// @Summary Получить список продукций
// @Description Возвращает список продукций с поддержкой пагинации
// @Tags products
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Success 200 {array} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, err := h.productUsecase.GetAllProducts(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Errorw("Ошибка получения списка продукций",
			"page", page,
			"limit", limit,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка продукций"})
		return
	}

	response := make([]ProductResponse, len(products))
	for i, product := range products {
		response[i] = h.productToResponse(product)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": response,
		"page":     page,
		"limit":    limit,
	})
}

// UpdateProduct обновляет продукцию
// @Summary Обновить продукцию
// @Description Обновляет информацию о продукции
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "ID продукции"
// @Param product body UpdateProductRequest true "Новые данные продукции"
// @Success 200 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Неверный ID продукции",
			"id", idStr,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукции"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Ошибка валидации запроса",
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	// Получаем существующую продукцию
	product, err := h.productUsecase.GetProductByID(c.Request.Context(), entity.ID(id))
	if err != nil {
		h.logger.Errorw("Ошибка получения продукции для обновления",
			"id", id,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукция не найдена"})
		return
	}

	// Обновляем поля продукции
	// Примечание: В реальном проекте нужно добавить методы для установки полей в entity
	// product.SetName(req.Name)
	// product.SetDescription(req.Description)
	// product.SetMinPartnerPrice(entity.Money(req.MinPartnerPrice))

	err = h.productUsecase.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		h.logger.Errorw("Ошибка обновления продукции",
			"id", id,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления продукции"})
		return
	}

	h.logger.Infow("Продукция успешно обновлена",
		"id", id,
		"request_id", c.GetString("request_id"),
	)

	response := h.productToResponse(product)
	c.JSON(http.StatusOK, response)
}

// DeleteProduct удаляет продукцию
// @Summary Удалить продукцию
// @Description Удаляет продукцию из системы
// @Tags products
// @Produce json
// @Param id path int true "ID продукции"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Неверный ID продукции",
			"id", idStr,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукции"})
		return
	}

	err = h.productUsecase.DeleteProduct(c.Request.Context(), entity.ID(id))
	if err != nil {
		h.logger.Errorw("Ошибка удаления продукции",
			"id", id,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления продукции"})
		return
	}

	h.logger.Infow("Продукция успешно удалена",
		"id", id,
		"request_id", c.GetString("request_id"),
	)

	c.Status(http.StatusNoContent)
}

// SearchProducts ищет продукции по запросу
// @Summary Поиск продукций
// @Description Поиск продукций по названию, артикулу или описанию
// @Tags products
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Success 200 {array} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
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

	products, err := h.productUsecase.SearchProducts(c.Request.Context(), query, limit, offset)
	if err != nil {
		h.logger.Errorw("Ошибка поиска продукций",
			"query", query,
			"page", page,
			"limit", limit,
			"error", err,
			"request_id", c.GetString("request_id"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска продукций"})
		return
	}

	response := make([]ProductResponse, len(products))
	for i, product := range products {
		response[i] = h.productToResponse(product)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": response,
		"query":    query,
		"page":     page,
		"limit":    limit,
	})
}

// productToResponse преобразует entity.Product в ProductResponse
func (h *ProductHandler) productToResponse(product *entity.Product) ProductResponse {
	response := ProductResponse{
		ID:              int(product.ID()),
		Article:         product.Article(),
		ProductTypeID:   int(product.ProductTypeID()),
		Name:            product.Name(),
		MinPartnerPrice: float64(product.MinPartnerPrice()),
		// Примечание: В реальном проекте нужно добавить методы для получения всех полей
		// Description:     product.Description(),
		// ImagePath:       product.ImagePath(),
		// WorkshopNumber:  product.WorkshopNumber(),
		// RequiredWorkers: product.RequiredWorkers(),
		// RollWidth:       product.RollWidth(),
		// CreatedAt:       product.CreatedAt().Format(time.RFC3339),
		// UpdatedAt:       product.UpdatedAt().Format(time.RFC3339),
	}

	return response
} 
