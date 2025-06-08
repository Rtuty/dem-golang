package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wallpaper-system/internal/config"
	"wallpaper-system/internal/database"
	"wallpaper-system/internal/handlers"
	"wallpaper-system/internal/repository"
	"wallpaper-system/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Wallpaper System API
// @version 1.0
// @description API для системы управления производством обоев "Наш декор"
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Устанавливаем режим Gin в зависимости от окружения
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Инициализируем логгер
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()

	// Устанавливаем глобальный логгер
	zap.ReplaceGlobals(logger)
	sugar := logger.Sugar()

	sugar.Infow("Запуск приложения",
		"name", "Wallpaper System",
		"version", "1.0.0",
		"environment", os.Getenv("APP_ENV"),
	)

	// Подключаемся к базе данных
	db, err := database.New(&cfg.Database)
	if err != nil {
		sugar.Fatalw("Ошибка подключения к базе данных", "error", err)
	}
	defer db.Close()

	// Инициализируем репозитории
	productRepo := repository.NewProductRepository(db.GetConnection())
	materialRepo := repository.NewMaterialRepository(db.GetConnection())

	// Инициализируем сервисы
	productService := services.NewProductService(productRepo, materialRepo)
	materialService := services.NewMaterialService(materialRepo)
	calculatorService := services.NewCalculatorService(materialRepo, productRepo)

	// Инициализируем хендлеры
	productHandler := handlers.NewProductHandler(productService, materialService)
	materialHandler := handlers.NewMaterialHandler(materialService, productService)
	calculatorHandler := handlers.NewCalculatorHandler(calculatorService, productService, materialService)

	// Создаем роутер Gin
	router := gin.Default()

	// Настраиваем CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Загружаем шаблоны
	router.LoadHTMLGlob("templates/*.html")

	// Подключаем статические файлы
	router.Static("/static", "./static")

	// Определяем функции для шаблонов
	router.SetFuncMap(map[string]interface{}{
		"add": func(a, b float64) float64 {
			return a + b
		},
	})

	// Настраиваем маршруты
	handlers.SetupRoutes(router, productHandler, materialHandler, calculatorHandler)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		sugar.Infow("Запуск HTTP сервера", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			sugar.Fatalw("Ошибка запуска HTTP сервера", "error", err)
		}
	}()

	// Выводим информацию о запуске
	printBanner(cfg)

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Получен сигнал завершения, останавливаем сервер...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalw("Ошибка при остановке сервера", "error", err)
	}

	sugar.Info("Сервер успешно остановлен")
}

// printBanner выводит баннер приложения
func printBanner(cfg *config.Config) {
	banner := fmt.Sprintf(`
╔══════════════════════════════════════════════════════════╗
║                  Wallpaper System                        ║
║                     Версия: 1.0.0                        ║
║                  Окружение: %s                      ║
╚══════════════════════════════════════════════════════════╝

🏭 Система управления производством обоев "Наш декор"
🌐 Адрес: http://%s:%s
📊 Доступные эндпоинты:
   • GET  /                          - Список продукции
   • GET  /products                  - Список продукции
   • GET  /products/:id              - Детали продукции
   • GET  /products/:id/materials    - Материалы для продукции
   • GET  /products/new              - Добавление продукции
   • GET  /products/:id/edit         - Редактирование продукции
   • GET  /calculator                - Калькулятор материалов

`,
		os.Getenv("APP_ENV"),
		cfg.Server.Host,
		cfg.Server.Port,
	)

	fmt.Print(banner)
}
