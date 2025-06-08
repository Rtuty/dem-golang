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

	"wallpaper-system/internal/infrastructure/config"
	"wallpaper-system/internal/infrastructure/logger"
	"wallpaper-system/internal/infrastructure/server"

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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем логгер
	zapLogger, err := logger.NewZapLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}

	defer zapLogger.Sync()

	// Устанавливаем глобальный логгер
	zap.ReplaceGlobals(zapLogger)
	sugaredLogger := zap.S()

	sugaredLogger.Infow("Запуск приложения",
		"name", cfg.App.Name,
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
	)

	// Инициализируем сервер с зависимостями
	srv, cleanup, err := server.NewServer(cfg, sugaredLogger)
	if err != nil {
		sugaredLogger.Fatalf("Ошибка инициализации сервера: %v", err)
	}
	defer cleanup()

	// Создаем HTTP сервер
	httpServer := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      srv,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		sugaredLogger.Infow("Запуск HTTP сервера",
			"address", httpServer.Addr,
		)

		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			sugaredLogger.Fatalw("Ошибка запуска HTTP сервера", "error", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	sugaredLogger.Info("Получен сигнал завершения, останавливаем сервер...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		sugaredLogger.Errorw("Ошибка при остановке сервера", "error", err)
		return
	}

	sugaredLogger.Info("Сервер успешно остановлен")
}

// setupEnvironment настраивает окружение приложения
func setupEnvironment(cfg *config.Config) {
	// Устанавливаем часовой пояс
	if cfg.App.Timezone != "" {
		if location, err := time.LoadLocation(cfg.App.Timezone); err == nil {
			time.Local = location
		}
	}

	// Устанавливаем переменные окружения для Gin
	if cfg.App.IsProduction() {
		os.Setenv("GIN_MODE", "release")
		return
	}

	if cfg.App.IsDevelopment() {
		os.Setenv("GIN_MODE", "debug")
	}
}

// printBanner выводит баннер приложения
func printBanner(cfg *config.Config) {
	banner := fmt.Sprintf(`
╔══════════════════════════════════════════════════════════╗
║                    %s                    ║
║                     Версия: %s                        ║
║                  Окружение: %s                ║
╚══════════════════════════════════════════════════════════╝

🏭 Система управления производством обоев "Наш декор"
🌐 Адрес: http://%s
📊 Доступные эндпоинты:
   • GET  /                   - Главная страница
   • GET  /api/v1/products    - API продукции
   • GET  /api/v1/materials   - API материалов
   • POST /api/v1/calculator  - API калькулятора

`,
		cfg.App.Name,
		cfg.App.Version,
		cfg.App.Environment,
		cfg.Server.GetServerAddress(),
	)

	fmt.Print(banner)
}
