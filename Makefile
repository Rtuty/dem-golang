# Переменные
APP_NAME := wallpaper-system
DOCKER_IMAGE := $(APP_NAME):latest
GO_VERSION := 1.21

# Цвета для вывода
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean dev-setup migrate-up migrate-down tools lint fmt

# Показать справку
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "  $(YELLOW)build$(NC)         - Собрать приложение локально"
	@echo "  $(YELLOW)run$(NC)           - Запустить приложение локально"
	@echo "  $(YELLOW)test$(NC)          - Запустить тесты"
	@echo "  $(YELLOW)clean$(NC)         - Очистить сборочные файлы"
	@echo ""
	@echo "  $(YELLOW)docker-build$(NC)  - Собрать Docker образ"
	@echo "  $(YELLOW)docker-run$(NC)    - Запустить с Docker Compose"
	@echo "  $(YELLOW)docker-stop$(NC)   - Остановить Docker контейнеры"
	@echo "  $(YELLOW)docker-clean$(NC)  - Очистить Docker ресурсы"
	@echo ""
	@echo "  $(YELLOW)dev-setup$(NC)     - Настроить среду разработки"
	@echo "  $(YELLOW)migrate-up$(NC)    - Выполнить миграции"
	@echo "  $(YELLOW)migrate-down$(NC)  - Откатить миграции"
	@echo "  $(YELLOW)tools$(NC)         - Запустить дополнительные инструменты"
	@echo ""
	@echo "  $(YELLOW)lint$(NC)          - Проверить код линтером"
	@echo "  $(YELLOW)fmt$(NC)           - Форматировать код"

# Локальная сборка
build:
	@echo "$(GREEN)Сборка приложения...$(NC)"
	go mod tidy
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate
	@echo "$(GREEN)Сборка завершена!$(NC)"

# Локальный запуск
run: build
	@echo "$(GREEN)Запуск приложения...$(NC)"
	./bin/server

# Запуск тестов
test:
	@echo "$(GREEN)Запуск тестов...$(NC)"
	go test -v ./...

# Очистка
clean:
	@echo "$(GREEN)Очистка сборочных файлов...$(NC)"
	rm -rf bin/
	go clean
	@echo "$(GREEN)Очистка завершена!$(NC)"

# Docker команды
docker-build:
	@echo "$(GREEN)Сборка Docker образа...$(NC)"
	docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Docker образ собран!$(NC)"

docker-run:
	@echo "$(GREEN)Запуск с Docker Compose...$(NC)"
	docker-compose up --build -d
	@echo "$(GREEN)Приложение доступно на http://localhost:8080$(NC)"

docker-stop:
	@echo "$(GREEN)Остановка Docker контейнеров...$(NC)"
	docker-compose down

docker-clean:
	@echo "$(GREEN)Очистка Docker ресурсов...$(NC)"
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Docker ресурсы очищены!$(NC)"

# Разработка
dev-setup:
	@echo "$(GREEN)Настройка среды разработки...$(NC)"
	go mod download
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Среда разработки настроена!$(NC)"

# Миграции
migrate-up:
	@echo "$(GREEN)Выполнение миграций...$(NC)"
	go run cmd/migrate/main.go up

migrate-down:
	@echo "$(GREEN)Откат миграций...$(NC)"
	go run cmd/migrate/main.go down

migrate-status:
	@echo "$(GREEN)Статус миграций:$(NC)"
	go run cmd/migrate/main.go status

# Дополнительные инструменты
tools:
	@echo "$(GREEN)Запуск дополнительных инструментов...$(NC)"
	docker-compose --profile tools up -d adminer
	@echo "$(GREEN)Adminer доступен на http://localhost:8081$(NC)"

# Линтер
lint:
	@echo "$(GREEN)Проверка кода линтером...$(NC)"
	golangci-lint run ./...

# Форматирование
fmt:
	@echo "$(GREEN)Форматирование кода...$(NC)"
	go fmt ./...
	goimports -w .

# Полный перезапуск
restart: docker-stop docker-run

# Логи Docker
logs:
	docker-compose logs -f app

# Логи базы данных
db-logs:
	docker-compose logs -f postgres

# Подключение к базе данных
db-connect:
	docker-compose exec postgres psql -U wallpaper_user -d wallpaper_system

# Бэкап базы данных
db-backup:
	@echo "$(GREEN)Создание бэкапа базы данных...$(NC)"
	docker-compose exec postgres pg_dump -U wallpaper_user wallpaper_system > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)Бэкап создан!$(NC)"

# Проверка состояния
status:
	@echo "$(GREEN)Состояние сервисов:$(NC)"
	docker-compose ps

# Быстрый старт
quick-start: docker-clean docker-run
	@echo "$(GREEN)Быстрый старт завершен!$(NC)"
	@echo "$(YELLOW)Приложение: http://localhost:8080$(NC)"
	@echo "$(YELLOW)База данных: localhost:5432$(NC)" 