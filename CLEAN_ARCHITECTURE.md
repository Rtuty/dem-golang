# Рефакторинг системы управления обоями в Clean Architecture

## Обзор изменений

Проект был полностью реорганизован для соответствия принципам Clean Architecture (Чистая архитектура) Роберта Мартина. Это обеспечивает:

- **Независимость от фреймворков** - бизнес-логика не зависит от внешних библиотек
- **Тестируемость** - легко тестировать бизнес-логику изолированно
- **Независимость от БД** - можно менять способ хранения данных
- **Независимость от UI** - можно менять интерфейс без изменения логики
- **Независимость от внешних агентов** - бизнес-правила не знают о внешнем мире

## Новая структура проекта

```
internal/
├── domain/                 # Слой доменов (внутренний круг)
│   ├── entities/          # Бизнес-сущности
│   │   ├── product.go     # Продукция с бизнес-логикой
│   │   ├── material.go    # Материалы с расчетами
│   │   └── errors.go      # Доменные ошибки
│   └── repositories/      # Интерфейсы репозиториев
│       ├── product_repository.go
│       └── material_repository.go
├── usecases/              # Слой вариантов использования
│   ├── product_usecase.go # Бизнес-логика продукции
│   └── material_usecase.go # Бизнес-логика материалов
├── adapters/              # Слой адаптеров интерфейсов
│   ├── controllers/       # HTTP контроллеры
│   │   ├── dto/          # Объекты передачи данных
│   │   ├── product_controller.go
│   │   └── calculator_controller.go
│   └── repositories/      # Реализации репозиториев
│       ├── product_repository_impl.go
│       └── material_repository_impl.go
└── infrastructure/        # Слой инфраструктуры (внешний круг)
    ├── config/           # Конфигурация
    ├── database/         # Подключение к БД
    └── server/           # HTTP сервер и маршруты
```

## Описание слоев

### 1. Domain Layer (Слой доменов)

**Назначение**: Содержит бизнес-сущности и правила предметной области.

**Файлы**:
- `entities/product.go` - Продукция с методами расчета цен и валидации
- `entities/material.go` - Материалы с расчетами необходимых количеств
- `entities/errors.go` - Доменные ошибки (валидация, бизнес-ошибки)
- `repositories/` - Интерфейсы для доступа к данным

**Особенности**:
- Независим от внешних слоев
- Содержит чистую бизнес-логику
- Определяет контракты (интерфейсы) для внешних слоев

### 2. Use Cases Layer (Слой вариантов использования)

**Назначение**: Содержит логику приложения и оркестрирует работу сущностей.

**Файлы**:
- `product_usecase.go` - Варианты использования для продукции
- `material_usecase.go` - Варианты использования для материалов и калькулятора

**Особенности**:
- Зависит только от Domain Layer
- Реализует конкретные сценарии использования
- Содержит валидацию и координацию между сущностями

### 3. Interface Adapters Layer (Слой адаптеров)

**Назначение**: Адаптирует данные между внешними интерфейсами и внутренними слоями.

**Компоненты**:
- **Controllers** - HTTP контроллеры для веб-интерфейса и API
- **DTOs** - Объекты передачи данных для преобразования между слоями
- **Repository Implementations** - Реализации интерфейсов репозиториев

**Особенности**:
- Преобразует данные между форматами
- Обрабатывает HTTP запросы и ответы
- Реализует интерфейсы, определенные в Domain Layer

### 4. Infrastructure Layer (Слой инфраструктуры)

**Назначение**: Содержит детали внешних систем и фреймворков.

**Компоненты**:
- **Config** - Конфигурация приложения
- **Database** - Подключение к базе данных
- **Server** - HTTP сервер и настройка маршрутов

**Особенности**:
- Самый внешний слой
- Содержит технические детали
- Может быть заменен без изменения бизнес-логики

## Принципы зависимостей

```
Infrastructure → Adapters → Use Cases → Domain
        ↓           ↓          ↓         ↓
      Веб,БД   Контроллеры  Бизнес-   Сущности
                Репозитории  логика
```

**Правило зависимостей**: Зависимости направлены только внутрь. Внутренние слои не знают о внешних.

## Преимущества новой архитектуры

### 1. Тестируемость
```go
// Легко тестировать бизнес-логику изолированно
func TestProductCalculatePrice(t *testing.T) {
    product := &entities.Product{
        Materials: []entities.ProductMaterial{...},
        ProductType: &entities.ProductType{Coefficient: 1.2},
    }
    
    price := product.CalculatePrice()
    assert.Equal(t, expectedPrice, price)
}
```

### 2. Гибкость
```go
// Можно легко заменить реализацию репозитория
type MockProductRepository struct{}

func (m *MockProductRepository) GetAll() ([]entities.Product, error) {
    return []entities.Product{...}, nil
}
```

### 3. Независимость от БД
Бизнес-логика работает с интерфейсами, а не с конкретными реализациями:

```go
type ProductUseCase struct {
    productRepo  repositories.ProductRepository  // Интерфейс!
    materialRepo repositories.MaterialRepository // Интерфейс!
}
```

### 4. Независимость от веб-фреймворка
Use Cases не знают о HTTP или Gin:

```go
func (uc *ProductUseCase) GetAllProducts() ([]entities.Product, error) {
    // Чистая бизнес-логика без HTTP
    return uc.productRepo.GetAll()
}
```

## Ключевые изменения

### 1. Разделение ответственности
- **Было**: Сервисы смешивали бизнес-логику с техническими деталями
- **Стало**: Четкое разделение по слоям с единственной ответственностью

### 2. Инверсия зависимостей
- **Было**: Прямые зависимости от конкретных реализаций
- **Стало**: Зависимости от абстракций (интерфейсов)

### 3. Доменная модель
- **Было**: Анемичные модели без поведения
- **Стало**: Богатые доменные модели с бизнес-логикой

### 4. Обработка ошибок
- **Было**: Технические ошибки смешивались с бизнес-ошибками
- **Стало**: Иерархия доменных ошибок с четкой семантикой

## Примеры использования

### Создание продукции
```go
// Use Case
func (uc *ProductUseCase) CreateProduct(product *entities.Product) error {
    // Доменная валидация
    if err := product.Validate(); err != nil {
        return err
    }
    
    // Бизнес-правила
    if _, err := uc.productRepo.GetProductTypeByID(product.ProductTypeID); err != nil {
        return entities.NewNotFoundError("тип продукции", strconv.Itoa(product.ProductTypeID))
    }
    
    return uc.productRepo.Create(product)
}

// Controller
func (c *ProductController) CreateProduct(ctx *gin.Context) {
    var dto dto.CreateProductRequestDTO
    if err := ctx.ShouldBindJSON(&dto); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
        return
    }
    
    product := dto.ToEntity()
    if err := c.productUseCase.CreateProduct(product); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusCreated, gin.H{"message": "Продукция создана"})
}
```

### Расчет материалов
```go
// Доменная логика в сущности
func (mcr *MaterialCalculationRequest) CalculateRequiredQuantity(
    productTypeCoeff, wastePercentage float64) int {
    
    materialPerUnit := mcr.ProductParam1 * mcr.ProductParam2 * productTypeCoeff
    totalMaterial := materialPerUnit * float64(mcr.ProductQuantity)
    materialWithWaste := totalMaterial * (1 + wastePercentage/100)
    requiredMaterial := materialWithWaste - mcr.MaterialInStock
    
    if requiredMaterial <= 0 {
        return 0
    }
    
    return int(requiredMaterial + 0.9999999)
}

// Use Case оркестрирует процесс
func (uc *CalculatorUseCase) CalculateRequiredMaterial(
    request *entities.MaterialCalculationRequest) (int, error) {
    
    productType, err := uc.materialRepo.GetProductTypeByID(request.ProductTypeID)
    if err != nil {
        return -1, err
    }
    
    materialType, err := uc.materialRepo.GetMaterialTypeByID(request.MaterialTypeID)
    if err != nil {
        return -1, err
    }
    
    return request.CalculateRequiredQuantity(
        productType.Coefficient,
        materialType.WastePercentage,
    ), nil
}
```

## Миграция

Проект был полностью рефакторован с сохранением всей функциональности:

1. ✅ Просмотр списка продукции с расчетом цен
2. ✅ Детальная информация о продукции
3. ✅ Калькулятор материалов с учетом брака
4. ✅ RESTful API для всех операций
5. ✅ Веб-интерфейс с формами и валидацией

Все существующие эндпоинты продолжают работать, но теперь построены на принципах чистой архитектуры.

## Запуск проекта

```bash
# Сборка
go build -o build/wallpaper-system.exe cmd/server/main.go

# Запуск
./build/wallpaper-system.exe

# Или напрямую
go run cmd/server/main.go
```

Приложение запустится с новым баннером, показывающим архитектурные слои и доступные эндпоинты. 