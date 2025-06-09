# 🏭 Наш декор - Система управления производством обоев

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Architecture](https://img.shields.io/badge/Architecture-Clean-green.svg)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![Version](https://img.shields.io/badge/Version-2.1-brightgreen.svg)](https://github.com/your-repo/wallpaper-system)

Современная система управления производством обоев, построенная на принципах Clean Architecture с оптимизированным фронтендом.

## ✨ Особенности

- 🏗️ **Clean Architecture** - четкое разделение слоев ответственности
- 📱 **Responsive UI** - адаптивный интерфейс для всех устройств  
- 🚀 **REST API** - полноценное API для интеграций
- 📊 **Управление продукцией** - создание, редактирование, калькулятор материалов
- 📦 **Складской учет** - контроль остатков и движения материалов
- 🐳 **Docker Ready** - готовые контейнеры для развертывания
- ⚡ **Оптимизированный код** - после рефакторинга v2.1

## 🏛️ Архитектура

```
dem-golang/
├── cmd/
│   ├── server/           # Точка входа приложения
│   └── migrate/          # Миграции БД
├── internal/
│   ├── adapters/         # Внешний слой (controllers, repositories)
│   ├── domain/          # Доменный слой (entities, interfaces)
│   ├── usecases/        # Бизнес-логика
│   └── infrastructure/  # Инфраструктура (DB, config, routes)
├── static/              # Статические файлы и SPA формы
├── templates/           # Server-side шаблоны
└── migrations/          # SQL миграции
```

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.21+
- Docker и Docker Compose
- PostgreSQL (в Docker)

### Установка и запуск

1. **Клонирование репозитория**
```bash
git clone <repository-url>
cd dem-golang
```

2. **Запуск через Docker Compose (Рекомендуется)**
```bash
# Запуск базы данных
docker-compose up -d postgres

# Выполнение миграций
docker-compose up -d migrate

# Запуск приложения
docker-compose up -d app
```

3. **Или локальный запуск**
```bash
# Настройка переменных окружения (опционально)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=wallpaper_user
export DB_PASSWORD=wallpaper_pass
export DB_NAME=wallpaper_system

# Компиляция и запуск
go build -o wallpaper-app cmd/server/main.go
./wallpaper-app
```

4. **Доступ к приложению**
- Основное приложение: http://localhost:8080
- Статические формы: http://localhost:8080/static/index.html
- API документация: http://localhost:8080/api/v1/

## 📋 API Эндпоинты

### 🌐 Веб-интерфейс
```
GET  /                     # Главная страница
GET  /products             # Список продукции
GET  /products/:id         # Детали продукции
GET  /materials            # Список материалов
GET  /calculator           # Калькулятор материалов
```

### 🔌 REST API
```
# Продукция
GET    /api/v1/products           # Список продукции
GET    /api/v1/products/:id       # Продукция по ID
POST   /api/v1/products           # Создать продукцию
PUT    /api/v1/products/:id       # Обновить продукцию
DELETE /api/v1/products/:id       # Удалить продукцию

# Материалы
GET    /api/v1/materials          # Список материалов
GET    /api/v1/materials/:id      # Материал по ID
POST   /api/v1/materials          # Создать материал
PUT    /api/v1/materials/:id      # Обновить материал
DELETE /api/v1/materials/:id      # Удалить материал

# Справочники
GET    /api/v1/product-types      # Типы продукции
GET    /api/v1/material-types     # Типы материалов
GET    /api/v1/measurement-units  # Единицы измерения
```

## 🎨 Фронтенд

Система включает два типа интерфейса:

### Server-Side Templates (основной)
- Полноценный веб-интерфейс с server-side рендерингом
- Адаптивный дизайн
- Интегрированный с бэкендом

### Static SPA Forms (дополнительный)
- Автономные HTML формы: `/static/index.html`
- Независимая работа с API
- Возможность работы без основного приложения

## 📊 База данных

### Структура
- **PostgreSQL** как основная СУБД
- **Миграции** для версионирования схемы
- **3NF нормализация** для целостности данных

### Основные таблицы
- `products` - Продукция
- `materials` - Материалы
- `product_types` - Типы продукции  
- `material_types` - Типы материалов
- `measurement_units` - Единицы измерения
- `product_materials` - Связи продукции с материалами

## 🔧 Конфигурация

Настройка через переменные окружения:

```bash
# Сервер
SERVER_HOST=localhost
SERVER_PORT=8080

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=wallpaper_user
DB_PASSWORD=wallpaper_pass
DB_NAME=wallpaper_system
DB_SSLMODE=disable
```

## 🏗️ Разработка

### Структура проекта
- **Чистая архитектура** - независимость слоев
- **Dependency Injection** - слабая связанность
- **Repository Pattern** - абстракция доступа к данным
- **Use Cases** - инкапсуляция бизнес-логики

### Добавление новой функциональности
1. Создать Entity в `domain/entities`
2. Добавить Repository interface в `domain/repositories`
3. Реализовать Repository в `adapters/repositories`
4. Создать Use Case в `usecases`
5. Добавить Controller в `adapters/controllers`
6. Обновить маршруты в `infrastructure/server`

## 📈 Версии

### v2.1 (Текущая) - Оптимизированный рефакторинг
- ✅ Удален дублированный код
- ✅ Объединены JavaScript файлы  
- ✅ Удалены неиспользуемые шаблоны
- ✅ Оптимизирован фронтенд
- ✅ Улучшена навигация
- ✅ Добавлен мониторинг статуса сервера

### v2.0 - Clean Architecture
- ✅ Полная реализация Clean Architecture
- ✅ REST API с полным CRUD
- ✅ Docker поддержка
- ✅ Статические формы
- ✅ Адаптивный дизайн

## 📝 Лицензия

© 2024 Наш декор. Все права защищены.

## 🤝 Вклад в проект

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit изменения (`git commit -m 'Add some AmazingFeature'`)
4. Push в branch (`git push origin feature/AmazingFeature`)
5. Открите Pull Request

## 📞 Поддержка

При возникновении вопросов или проблем создайте Issue в репозитории. 