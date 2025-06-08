package server

import (
	"wallpaper-system/internal/adapters/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты приложения
func SetupRoutes(
	router *gin.Engine,
	productController *controllers.ProductController,
	calculatorController *controllers.CalculatorController,
) {
	// Главная страница - перенаправление на продукцию
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/products")
	})

	// Веб-страницы
	setupWebRoutes(router, productController, calculatorController)

	// API маршруты
	setupAPIRoutes(router, productController, calculatorController)
}

// setupWebRoutes настраивает веб-маршруты
func setupWebRoutes(
	router *gin.Engine,
	productController *controllers.ProductController,
	calculatorController *controllers.CalculatorController,
) {
	// Продукция
	router.GET("/products", productController.GetProductsPage)
	router.GET("/products/:id", productController.GetProductDetailsPage)

	// Калькулятор
	router.GET("/calculator", calculatorController.GetCalculatorPage)
	router.POST("/calculator", calculatorController.CalculateMaterial)
}

// setupAPIRoutes настраивает API маршруты
func setupAPIRoutes(
	router *gin.Engine,
	productController *controllers.ProductController,
	calculatorController *controllers.CalculatorController,
) {
	api := router.Group("/api/v1")
	{
		// Продукция API
		products := api.Group("/products")
		{
			products.GET("", productController.GetProducts)
			products.GET("/:id", productController.GetProductByID)
			products.POST("", productController.CreateProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}

		// Калькулятор API
		calculator := api.Group("/calculator")
		{
			calculator.POST("/calculate", calculatorController.CalculateMaterialAPI)
		}
	}
}
