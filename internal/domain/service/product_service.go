package service

import (
	"context"
	"errors"
	"fmt"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"
)

// ProductService представляет доменный сервис для работы с продукцией
type ProductService struct {
	productRepo     repository.ProductRepository
	productTypeRepo repository.ProductTypeRepository
	materialRepo    repository.MaterialRepository
}

// NewProductService создает новый доменный сервис продукции
func NewProductService(
	productRepo repository.ProductRepository,
	productTypeRepo repository.ProductTypeRepository,
	materialRepo repository.MaterialRepository,
) *ProductService {
	return &ProductService{
		productRepo:     productRepo,
		productTypeRepo: productTypeRepo,
		materialRepo:    materialRepo,
	}
}

// CreateProduct создает новую продукцию
func (s *ProductService) CreateProduct(
	ctx context.Context,
	article string,
	productTypeID entity.ID,
	name string,
	minPartnerPrice entity.Money,
	description *string,
) (*entity.Product, error) {
	// Проверяем уникальность артикула
	exists, err := s.productRepo.ExistsByArticle(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки уникальности артикула: %w", err)
	}
	if exists {
		return nil, errors.New("продукция с таким артикулом уже существует")
	}

	// Проверяем существование типа продукции
	productType, err := s.productTypeRepo.GetByID(ctx, productTypeID)
	if err != nil {
		return nil, fmt.Errorf("тип продукции не найден: %w", err)
	}

	// Создаем новую продукцию
	product, err := entity.NewProduct(article, productTypeID, name, minPartnerPrice)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания продукции: %w", err)
	}

	if description != nil {
		if err := product.UpdateBasicInfo(name, description); err != nil {
			return nil, fmt.Errorf("ошибка установки описания: %w", err)
		}
	}

	product.SetProductType(productType)

	// Сохраняем в репозитории
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("ошибка сохранения продукции: %w", err)
	}

	return product, nil
}

// GetProductByID возвращает продукцию по ID
func (s *ProductService) GetProductByID(ctx context.Context, id entity.ID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продукции: %w", err)
	}

	// Загружаем материалы и рассчитываем цену
	materials, err := s.productRepo.GetMaterials(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов: %w", err)
	}

	product.SetMaterials(materials)

	// Рассчитываем цену, если есть материалы
	if len(materials) > 0 {
		if err := product.CalculatePrice(); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			// log.Printf("Не удалось рассчитать цену для продукции %d: %v", id, err)
		}
	}

	return product, nil
}

// GetAllProducts возвращает все продукции с пагинацией
func (s *ProductService) GetAllProducts(ctx context.Context, limit, offset int) ([]*entity.Product, int, error) {
	products, err := s.productRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения списка продукции: %w", err)
	}

	// Для каждого продукта загружаем материалы и рассчитываем цену
	for _, product := range products {
		materials, err := s.productRepo.GetMaterials(ctx, product.ID())
		if err == nil && len(materials) > 0 {
			product.SetMaterials(materials)
			product.CalculatePrice() // Игнорируем ошибку расчета
		}
	}

	// Получаем общее количество
	total, err := s.productRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения общего количества продукции: %w", err)
	}

	return products, total, nil
}

// UpdateProduct обновляет продукцию
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	id entity.ID,
	name *string,
	description *string,
	minPartnerPrice *entity.Money,
) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	// Обновляем основную информацию
	if name != nil {
		if err := product.UpdateBasicInfo(*name, description); err != nil {
			return fmt.Errorf("ошибка обновления основной информации: %w", err)
		}
	}

	// Обновляем цену
	if minPartnerPrice != nil {
		if err := product.UpdatePrice(*minPartnerPrice); err != nil {
			return fmt.Errorf("ошибка обновления цены: %w", err)
		}
	}

	// Сохраняем изменения
	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}

	return nil
}

// DeleteProduct удаляет продукцию
func (s *ProductService) DeleteProduct(ctx context.Context, id entity.ID) error {
	// Проверяем существование продукции
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	// Удаляем продукцию
	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ошибка удаления продукции: %w", err)
	}

	return nil
}

// AddMaterialToProduct добавляет материал к продукции
func (s *ProductService) AddMaterialToProduct(
	ctx context.Context,
	productID, materialID entity.ID,
	quantityPerUnit float64,
) error {
	// Проверяем существование продукции
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	// Проверяем существование материала
	_, err = s.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	// Создаем связь продукции с материалом
	productMaterial, err := entity.NewProductMaterial(productID, materialID, quantityPerUnit)
	if err != nil {
		return fmt.Errorf("ошибка создания связи: %w", err)
	}

	// Сохраняем связь
	if err := s.productRepo.AddMaterial(ctx, productMaterial); err != nil {
		return fmt.Errorf("ошибка добавления материала к продукции: %w", err)
	}

	return nil
}

// RemoveMaterialFromProduct удаляет материал из продукции
func (s *ProductService) RemoveMaterialFromProduct(
	ctx context.Context,
	productID, materialID entity.ID,
) error {
	// Проверяем существование продукции
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	// Удаляем связь
	if err := s.productRepo.RemoveMaterial(ctx, productID, materialID); err != nil {
		return fmt.Errorf("ошибка удаления материала из продукции: %w", err)
	}

	return nil
}

// UpdateMaterialQuantityInProduct обновляет количество материала в продукции
func (s *ProductService) UpdateMaterialQuantityInProduct(
	ctx context.Context,
	productID, materialID entity.ID,
	newQuantity float64,
) error {
	if newQuantity <= 0 {
		return errors.New("количество материала должно быть положительным")
	}

	// Проверяем существование продукции
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("продукция не найдена: %w", err)
	}

	// Обновляем количество
	if err := s.productRepo.UpdateMaterialQuantity(ctx, productID, materialID, newQuantity); err != nil {
		return fmt.Errorf("ошибка обновления количества материала: %w", err)
	}

	return nil
}

// CalculateProductPrice рассчитывает цену продукции на основе материалов
func (s *ProductService) CalculateProductPrice(ctx context.Context, productID entity.ID) (entity.Money, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("продукция не найдена: %w", err)
	}

	// Загружаем материалы
	materials, err := s.productRepo.GetMaterials(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения материалов: %w", err)
	}

	product.SetMaterials(materials)

	// Рассчитываем цену
	if err := product.CalculatePrice(); err != nil {
		return 0, fmt.Errorf("ошибка расчета цены: %w", err)
	}

	if product.CalculatedPrice() == nil {
		return 0, errors.New("цена не была рассчитана")
	}

	return *product.CalculatedPrice(), nil
} 