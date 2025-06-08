package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"wallpaper-system/internal/application/usecase"
	"wallpaper-system/internal/domain/service"
	"wallpaper-system/internal/infrastructure/config"
	"wallpaper-system/internal/infrastructure/database"
	"wallpaper-system/internal/infrastructure/repository"
	transportHTTP "wallpaper-system/internal/infrastructure/transport/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server представляет HTTP сервер
type Server struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
	router *gin.Engine
	db     *sql.DB
}

// NewServer создает новый HTTP сервер с dependency injection
func NewServer(cfg *config.Config, logger *zap.SugaredLogger) (*gin.Engine, func(), error) {
	// Подключаемся к базе данных
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Инициализируем репозитории через реестр
	repoRegistry := repository.NewRepositoryRegistry(db)

	// Инициализируем доменные сервисы
	productService := service.NewProductService(repoRegistry.Product)
	materialService := service.NewMaterialService(repoRegistry.Material)
	calculatorService := service.NewCalculatorService()

	// Инициализируем use cases
	productUseCase := usecase.NewProductUsecase(productService, calculatorService)
	calculatorUseCase := usecase.NewCalculatorUsecase(calculatorService)

	// Настраиваем Gin в зависимости от окружения
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.App.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	}

	// Создаем роутер
	router := gin.New()

	// Настраиваем middleware
	setupMiddleware(router, cfg, logger.Desugar())

	// Настраиваем маршруты
	setupRoutes(router, productUseCase, calculatorUseCase, logger)

	// Функция очистки ресурсов
	cleanup := func() {
		if err := db.Close(); err != nil {
			logger.Errorw("Ошибка закрытия соединения с БД", "error", err)
		}
	}

	return router, cleanup, nil
}

// setupMiddleware настраивает middleware для роутера
func setupMiddleware(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// Recovery middleware
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered", zap.String("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": "Произошла внутренняя ошибка сервера",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// Logger middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	}))

	// CORS middleware
	if cfg.Server.EnableCORS {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"*"} // В продакшене указать конкретные домены
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
		router.Use(cors.New(corsConfig))
	}

	// Timeout middleware
	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	// Request ID middleware
	router.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	})

	// Security headers middleware
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})
}

// setupRoutes настраивает маршруты
func setupRoutes(
	router *gin.Engine,
	productUseCase *usecase.ProductUseCase,
	calculatorUseCase *usecase.CalculatorUseCase,
	logger *zap.SugaredLogger,
) {
	// Статические файлы
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*.html")

	// Создаем обработчики
	productHandler := transportHTTP.NewProductHandler(productUseCase, logger)
	calculatorHandler := transportHTTP.NewCalculatorHandler(calculatorUseCase, logger)
	healthHandler := transportHTTP.NewHealthHandler(logger)

	// Главная страница и веб-интерфейс
	setupWebRoutes(router, productHandler, calculatorHandler)

	// API маршруты
	setupAPIRoutes(router, productHandler, calculatorHandler, healthHandler)
}

// setupWebRoutes настраивает веб-маршруты
func setupWebRoutes(
	router *gin.Engine,
	productHandler *transportHTTP.ProductHandler,
	calculatorHandler *transportHTTP.CalculatorHandler,
) {
	// Главная страница
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Система управления производством обоев",
		})
	})

	// Веб-интерфейс для продукции
	web := router.Group("/web")
	{
		web.GET("/products", productHandler.GetProductsPage)
		web.GET("/products/new", productHandler.CreateProductPage)
		web.POST("/products", productHandler.CreateProductWeb)
		web.GET("/products/:id", productHandler.GetProductPage)
		web.GET("/products/:id/edit", productHandler.EditProductPage)
		web.PUT("/products/:id", productHandler.UpdateProductWeb)
		web.DELETE("/products/:id", productHandler.DeleteProductWeb)
	}

	// Веб-интерфейс для калькулятора
	webCalc := router.Group("/web/calculator")
	{
		webCalc.GET("/", calculatorHandler.CalculatorPage)
		webCalc.POST("/calculate", calculatorHandler.CalculateWeb)
	}
}

// setupAPIRoutes настраивает API маршруты
func setupAPIRoutes(
	router *gin.Engine,
	productHandler *transportHTTP.ProductHandler,
	calculatorHandler *transportHTTP.CalculatorHandler,
	healthHandler *transportHTTP.HealthHandler,
) {
	// API v1
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthHandler.Health)
		api.GET("/ready", healthHandler.Ready)

		// Продукция
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
			products.POST("/:id/materials", productHandler.AddMaterialToProduct)
			products.DELETE("/:id/materials/:material_id", productHandler.RemoveMaterialFromProduct)
			products.GET("/:id/price", productHandler.CalculateProductPrice)
		}

		// Калькулятор
		calculator := api.Group("/calculator")
		{
			calculator.POST("/materials", calculatorHandler.CalculateRequiredMaterial)
			calculator.GET("/materials/summary", calculatorHandler.GetCalculationSummary)
		}

		// Метрики и информация о системе
		system := api.Group("/system")
		{
			system.GET("/info", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"name":        "Wallpaper System",
					"version":     "1.0.0",
					"environment": "development",
					"build_time":  time.Now().Format(time.RFC3339),
				})
			})
		}
	}

	// Swagger документация (только в режиме разработки)
	if gin.Mode() != gin.ReleaseMode {
		// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ErrorResponse представляет стандартный формат ошибки
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// SuccessResponse представляет стандартный формат успешного ответа
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginationMeta представляет метаданные пагинации
type PaginationMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
} 