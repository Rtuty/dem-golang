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

	// Слой инфраструктуры
	"wallpaper-system/internal/infrastructure/config"
	"wallpaper-system/internal/infrastructure/database"
	"wallpaper-system/internal/infrastructure/server"

	// Слой адаптеров
	"wallpaper-system/internal/adapters/controllers"
	"wallpaper-system/internal/adapters/repositories"

	// Слой вариантов использования
	"wallpaper-system/internal/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Wallpaper System API
// @version 2.0
// @description API для системы управления производством обоев "Наш декор" с чистой архитектурой
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

	sugar.Infow("Запуск приложения с чистой архитектурой",
		"name", "Wallpaper System",
		"version", "2.0.0",
		"environment", os.Getenv("APP_ENV"),
	)

	// Подключаемся к базе данных (слой инфраструктуры)
	db, err := database.New(&cfg.Database)
	if err != nil {
		sugar.Fatalw("Ошибка подключения к базе данных", "error", err)
	}
	defer db.Close()

	// Инициализируем репозитории (слой адаптеров)
	productRepo := repositories.NewProductRepository(db.GetConnection())
	materialRepo := repositories.NewMaterialRepository(db.GetConnection())

	// Инициализируем варианты использования (слой бизнес-логики)
	productUseCase := usecases.NewProductUseCase(productRepo, materialRepo)
	materialUseCase := usecases.NewMaterialUseCase(materialRepo)
	calculatorUseCase := usecases.NewCalculatorUseCase(materialRepo)

	// Инициализируем контроллеры (слой адаптеров)
	productController := controllers.NewProductController(productUseCase, materialUseCase)
	materialController := controllers.NewMaterialController(materialUseCase)
	calculatorController := controllers.NewCalculatorController(calculatorUseCase, materialUseCase, productUseCase)

	// Создаем роутер Gin
	router := gin.Default()

	// Настраиваем CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Определяем функции для шаблонов (ДО загрузки шаблонов!)
	router.SetFuncMap(map[string]interface{}{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
	})

	// Загружаем шаблоны
	router.LoadHTMLGlob("templates/*.html")

	// Подключаем статические файлы
	router.Static("/static", "./static")

	// Настраиваем маршруты (слой инфраструктуры)
	server.SetupRoutes(router, productController, calculatorController, materialController)

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
║               Wallpaper System v2.0                      ║
║              Чистая архитектура                           ║
║                  Окружение: %s                      ║
╚══════════════════════════════════════════════════════════╝

🏭 Система управления производством обоев "Наш декор"
🌐 Адрес: http://%s:%s
📊 Архитектурные слои:
   • Domain Layer        - Бизнес-сущности и правила
   • Use Cases Layer     - Варианты использования
   • Interface Adapters  - Контроллеры и репозитории  
   • Infrastructure      - Веб, БД, конфигурация

📋 Доступные эндпоинты:
   • GET  /                          - Главная страница
   • GET  /products                  - Список продукции
   • GET  /products/:id              - Детали продукции
   • GET  /calculator                - Калькулятор материалов
   • POST /calculator                - Расчет материалов
   • API  /api/v1/products           - REST API продукции
   • API  /api/v1/calculator         - REST API калькулятора

`,
		os.Getenv("APP_ENV"),
		cfg.Server.Host,
		cfg.Server.Port,
	)

	fmt.Print(banner)
}
