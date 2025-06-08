package main

import (
	"html/template"
	"log"

	"wallpaper-system/internal/config"
	"wallpaper-system/internal/database"
	"wallpaper-system/internal/handlers"
	"wallpaper-system/internal/repository"
	"wallpaper-system/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к базе данных
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	productRepo := repository.NewProductRepository(db.GetConnection())
	materialRepo := repository.NewMaterialRepository(db.GetConnection())

	// Инициализация сервисов
	productService := services.NewProductService(productRepo, materialRepo)
	materialService := services.NewMaterialService(materialRepo)

	// Инициализация хендлеров
	productHandler := handlers.NewProductHandler(productService, materialService)
	materialHandler := handlers.NewMaterialHandler(materialService, productService)

	// Настройка Gin
	if gin.Mode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Загрузка HTML шаблонов с функциями
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b float64) float64 {
			return a + b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	})
	router.LoadHTMLGlob("templates/*.html")

	// Статические файлы
	router.Static("/static", "./static")

	// Маршруты для веб-интерфейса
	setupWebRoutes(router, productHandler, materialHandler)

	// API маршруты
	setupAPIRoutes(router, productHandler, materialHandler, productService)

	// Запуск сервера
	serverAddr := cfg.Server.GetServerAddress()
	log.Printf("Сервер запущен на %s", serverAddr)
	log.Printf("Откройте браузер и перейдите по адресу: http://%s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// setupWebRoutes настраивает маршруты для веб-интерфейса
func setupWebRoutes(router *gin.Engine, productHandler *handlers.ProductHandler, materialHandler *handlers.MaterialHandler) {
	// Главная страница - список продукции
	router.GET("/", productHandler.GetProducts)

	// Продукция
	router.GET("/products", productHandler.GetProducts)
	router.GET("/products/new", productHandler.ShowCreateProductForm)
	router.POST("/products/new", productHandler.CreateProduct)
	router.GET("/products/:id", productHandler.GetProduct)
	router.GET("/products/:id/edit", productHandler.ShowEditProductForm)
	router.POST("/products/:id/edit", productHandler.UpdateProduct)
	router.GET("/products/:id/materials", productHandler.GetProductMaterials)

	// Материалы
	router.GET("/materials", materialHandler.GetMaterials)

	// Калькулятор материалов
	router.GET("/calculator", materialHandler.ShowCalculatorForm)
	router.POST("/calculator", materialHandler.CalculateMaterialForm)
}

// setupAPIRoutes настраивает API маршруты
func setupAPIRoutes(router *gin.Engine, productHandler *handlers.ProductHandler, materialHandler *handlers.MaterialHandler, productService *services.ProductService) {
	api := router.Group("/api")

	// API для продукции
	products := api.Group("/products")
	{
		products.GET("/", func(c *gin.Context) {
			// Возвращаем JSON список продукции
			productsList, err := productService.GetAllProducts()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"products": productsList})
		})
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	// API для калькулятора материалов
	calculator := api.Group("/calculator")
	{
		calculator.POST("/", materialHandler.CalculateMaterial)
	}
}
