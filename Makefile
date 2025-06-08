# Wallpaper System Makefile
# Поддержка Windows, Linux и macOS

.PHONY: help build run test clean install migrate dev docker-build docker-run docker-stop fmt lint tidy

# Определяем операционную систему
UNAME_S := $(shell uname -s 2>nul || echo Windows)
ifeq ($(UNAME_S),Windows)
	EXE := .exe
	RM := del /Q
	RMDIR := rmdir /S /Q
	MKDIR := mkdir
	NULL := nul
else
	EXE :=
	RM := rm -f
	RMDIR := rm -rf
	MKDIR := mkdir -p
	NULL := /dev/null
endif

# Переменные
APP_NAME := wallpaper-system
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)$(EXE)
MAIN_PATH := ./cmd/server
MIGRATE_PATH := ./cmd/migrate

# Go параметры
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# Версия и билд информация
VERSION := $(shell git describe --tags --always --dirty 2>$(NULL) || echo "unknown")
BUILD_TIME := $(shell date +%Y%m%d-%H%M%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>$(NULL) || echo "unknown")

# LDFLAGS для встраивания информации о сборке
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# По умолчанию показываем справку
help: ## Показать справку
	@echo "=== Wallpaper System Build Commands ==="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Примеры использования:"
	@echo "  make install      # Установить зависимости"
	@echo "  make dev          # Запустить в режиме разработки"
	@echo "  make build        # Собрать приложение"
	@echo "  make test         # Запустить тесты"
	@echo "  make docker-run   # Запустить с Docker"

install: ## Установить зависимости
	@echo "Установка зависимостей..."
	$(GOGET) -v ./...
	$(GOMOD) tidy
	$(GOMOD) download
	@echo "Зависимости установлены успешно!"

build: ## Собрать приложение
	@echo "Сборка приложения..."
	@if not exist $(BUILD_DIR) $(MKDIR) $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY) $(MAIN_PATH)
	@echo "Приложение собрано: $(BINARY)"

build-all: ## Собрать для всех платформ
	@echo "Сборка для всех платформ..."
	@if not exist $(BUILD_DIR) $(MKDIR) $(BUILD_DIR)
	
	@echo "Сборка для Windows..."
	SET GOOS=windows&& SET GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "Сборка для Linux..."
	SET GOOS=linux&& SET GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)
	
	@echo "Сборка для macOS..."
	SET GOOS=darwin&& SET GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PATH)
	
	@echo "Сборка завершена!"

run: build ## Собрать и запустить приложение
	@echo "Запуск приложения..."
	$(BINARY)

dev: ## Запустить в режиме разработки (без сборки)
	@echo "Запуск в режиме разработки..."
	SET APP_ENV=development&& SET LOG_LEVEL=debug&& $(GORUN) $(MAIN_PATH)

test: ## Запустить тесты
	@echo "Запуск тестов..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Тесты завершены!"

test-coverage: test ## Запустить тесты с отчетом о покрытии
	@echo "Генерация отчета о покрытии..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Отчет сохранен в coverage.html"

bench: ## Запустить бенчмарки
	@echo "Запуск бенчмарков..."
	$(GOTEST) -bench=. -benchmem ./...

migrate-up: ## Применить миграции БД
	@echo "Применение миграций..."
	$(GORUN) $(MIGRATE_PATH) up

migrate-down: ## Откатить миграции БД
	@echo "Откат миграций..."
	$(GORUN) $(MIGRATE_PATH) down

migrate-create: ## Создать новую миграцию (make migrate-create NAME=add_users_table)
	@echo "Создание миграции: $(NAME)"
	$(GORUN) $(MIGRATE_PATH) create $(NAME)

clean: ## Очистить сборочные файлы
	@echo "Очистка..."
	$(GOCLEAN)
	@if exist $(BUILD_DIR) $(RMDIR) $(BUILD_DIR)
	@if exist coverage.out $(RM) coverage.out
	@if exist coverage.html $(RM) coverage.html
	@echo "Очистка завершена!"

fmt: ## Форматировать код
	@echo "Форматирование кода..."
	$(GOFMT) ./...
	@echo "Форматирование завершено!"

lint: ## Запустить линтер (требует установки golangci-lint)
	@echo "Запуск линтера..."
	@golangci-lint run || echo "Установите golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

tidy: ## Очистить go.mod
	@echo "Очистка зависимостей..."
	$(GOMOD) tidy
	@echo "Зависимости очищены!"

update: ## Обновить зависимости
	@echo "Обновление зависимостей..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "Зависимости обновлены!"

docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	docker build -t $(APP_NAME):latest .
	@echo "Docker образ собран!"

docker-run: docker-build ## Запустить с Docker Compose
	@echo "Запуск с Docker Compose..."
	docker-compose up --build

docker-stop: ## Остановить Docker Compose
	@echo "Остановка Docker Compose..."
	docker-compose down

docker-clean: ## Очистить Docker ресурсы
	@echo "Очистка Docker ресурсов..."
	docker-compose down -v --remove-orphans
	docker image prune -f

# Задачи для разработки
setup: install migrate-up ## Полная настройка проекта для разработки
	@echo "Проект настроен для разработки!"

dev-deps: ## Установить зависимости для разработки
	@echo "Установка инструментов разработки..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/swaggo/swag/cmd/swag
	@echo "Инструменты разработки установлены!"

swagger: ## Генерация Swagger документации
	@echo "Генерация Swagger документации..."
	@swag init -g cmd/server/main.go -o docs || echo "Установите swag: go install github.com/swaggo/swag/cmd/swag@latest"

# Задачи для продакшена
prod-build: ## Сборка для продакшена
	@echo "Сборка для продакшена..."
	@if not exist $(BUILD_DIR) $(MKDIR) $(BUILD_DIR)
	SET CGO_ENABLED=0&& SET GOOS=linux&& $(GOBUILD) -a -installsuffix cgo $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-prod $(MAIN_PATH)
	@echo "Продакшен сборка готова!"

# Информация о проекте
info: ## Показать информацию о проекте
	@echo "=== Информация о проекте ==="
	@echo "Название: $(APP_NAME)"
	@echo "Версия: $(VERSION)"
	@echo "Коммит: $(GIT_COMMIT)"
	@echo "Время сборки: $(BUILD_TIME)"
	@echo "Go версия: $(shell $(GOCMD) version)"
	@echo "Платформа: $(UNAME_S)"

# Задачи для CI/CD
ci: fmt lint test ## Задачи для CI (форматирование, линтинг, тесты)
	@echo "CI задачи выполнены успешно!"

# Алиасы для удобства
run-dev: dev ## Алиас для dev
start: run ## Алиас для run
stop: docker-stop ## Алиас для docker-stop 