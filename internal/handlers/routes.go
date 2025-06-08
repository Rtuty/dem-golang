package handlers

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты приложения
func SetupRoutes(router *gin.Engine, productHandler *ProductHandler, materialHandler *MaterialHandler, calculatorHandler *CalculatorHandler) {
	// Главная страница (список продукции)
	router.GET("/", productHandler.GetProducts)

	// Маршруты для продукции
	products := router.Group("/products")
	{
		products.GET("/", productHandler.GetProducts)
		products.GET("/new", productHandler.ShowCreateProductForm)
		products.POST("/", productHandler.CreateProduct)
		products.GET("/:id", productHandler.GetProduct)
		products.GET("/:id/edit", productHandler.ShowEditProductForm)
		products.POST("/:id", productHandler.UpdateProduct)
		products.GET("/:id/materials", productHandler.GetProductMaterials)

		// API для AJAX запросов
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Маршруты для материалов
	materials := router.Group("/materials")
	{
		materials.GET("/", materialHandler.GetMaterials)
	}

	// Маршруты для калькулятора
	calculator := router.Group("/calculator")
	{
		calculator.GET("/", calculatorHandler.ShowCalculatorForm)
		calculator.POST("/calculate", calculatorHandler.CalculateMaterialForm)

		// API для AJAX запросов
		calculator.POST("/api/calculate", calculatorHandler.CalculateMaterialAmount)
	}
}
