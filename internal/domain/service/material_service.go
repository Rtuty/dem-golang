package service

import (
	"context"
	"errors"
	"fmt"

	"wallpaper-system/internal/domain/entity"
	"wallpaper-system/internal/domain/repository"
)

// MaterialService представляет доменный сервис для работы с материалами
type MaterialService struct {
	materialRepo     repository.MaterialRepository
	materialTypeRepo repository.MaterialTypeRepository
	measurementRepo  repository.MeasurementUnitRepository
}

// NewMaterialService создает новый доменный сервис материалов
func NewMaterialService(
	materialRepo repository.MaterialRepository,
	materialTypeRepo repository.MaterialTypeRepository,
	measurementRepo repository.MeasurementUnitRepository,
) *MaterialService {
	return &MaterialService{
		materialRepo:     materialRepo,
		materialTypeRepo: materialTypeRepo,
		measurementRepo:  measurementRepo,
	}
}

// CreateMaterial создает новый материал
func (s *MaterialService) CreateMaterial(
	ctx context.Context,
	article string,
	materialTypeID entity.ID,
	name string,
	measurementUnitID entity.ID,
	packageQuantity float64,
	costPerUnit entity.Money,
	description *string,
) (*entity.Material, error) {
	// Проверяем уникальность артикула
	exists, err := s.materialRepo.ExistsByArticle(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки уникальности артикула: %w", err)
	}
	if exists {
		return nil, errors.New("материал с таким артикулом уже существует")
	}

	// Проверяем существование типа материала
	materialType, err := s.materialTypeRepo.GetByID(ctx, materialTypeID)
	if err != nil {
		return nil, fmt.Errorf("тип материала не найден: %w", err)
	}

	// Проверяем существование единицы измерения
	measurementUnit, err := s.measurementRepo.GetByID(ctx, measurementUnitID)
	if err != nil {
		return nil, fmt.Errorf("единица измерения не найдена: %w", err)
	}

	// Создаем новый материал
	material, err := entity.NewMaterial(
		article,
		materialTypeID,
		name,
		measurementUnitID,
		packageQuantity,
		costPerUnit,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания материала: %w", err)
	}

	if description != nil {
		if err := material.UpdateBasicInfo(name, description); err != nil {
			return nil, fmt.Errorf("ошибка установки описания: %w", err)
		}
	}

	material.SetMaterialType(materialType)
	material.SetMeasurementUnit(measurementUnit)

	// Сохраняем в репозитории
	if err := s.materialRepo.Create(ctx, material); err != nil {
		return nil, fmt.Errorf("ошибка сохранения материала: %w", err)
	}

	return material, nil
}

// GetMaterialByID возвращает материал по ID
func (s *MaterialService) GetMaterialByID(ctx context.Context, id entity.ID) (*entity.Material, error) {
	material, err := s.materialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материала: %w", err)
	}

	return material, nil
}

// GetAllMaterials возвращает все материалы с пагинацией
func (s *MaterialService) GetAllMaterials(ctx context.Context, limit, offset int) ([]*entity.Material, int, error) {
	materials, err := s.materialRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения списка материалов: %w", err)
	}

	// Получаем общее количество
	total, err := s.materialRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения общего количества материалов: %w", err)
	}

	return materials, total, nil
}

// UpdateMaterial обновляет материал
func (s *MaterialService) UpdateMaterial(
	ctx context.Context,
	id entity.ID,
	name *string,
	description *string,
	costPerUnit *entity.Money,
) error {
	material, err := s.materialRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	// Обновляем основную информацию
	if name != nil {
		if err := material.UpdateBasicInfo(*name, description); err != nil {
			return fmt.Errorf("ошибка обновления основной информации: %w", err)
		}
	}

	// Обновляем стоимость
	if costPerUnit != nil {
		if err := material.UpdateCost(*costPerUnit); err != nil {
			return fmt.Errorf("ошибка обновления стоимости: %w", err)
		}
	}

	// Сохраняем изменения
	if err := s.materialRepo.Update(ctx, material); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}

	return nil
}

// DeleteMaterial удаляет материал
func (s *MaterialService) DeleteMaterial(ctx context.Context, id entity.ID) error {
	// Проверяем существование материала
	_, err := s.materialRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	// Удаляем материал
	if err := s.materialRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ошибка удаления материала: %w", err)
	}

	return nil
}

// UpdateStock обновляет остаток материала на складе
func (s *MaterialService) UpdateStock(ctx context.Context, materialID entity.ID, newQuantity float64) error {
	material, err := s.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	if err := material.UpdateStock(newQuantity); err != nil {
		return fmt.Errorf("ошибка обновления остатка: %w", err)
	}

	if err := s.materialRepo.Update(ctx, material); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}

	return nil
}

// AddToStock добавляет материал на склад
func (s *MaterialService) AddToStock(ctx context.Context, materialID entity.ID, quantity float64) error {
	material, err := s.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	if err := material.AddToStock(quantity); err != nil {
		return fmt.Errorf("ошибка добавления на склад: %w", err)
	}

	if err := s.materialRepo.Update(ctx, material); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}

	return nil
}

// RemoveFromStock списывает материал со склада
func (s *MaterialService) RemoveFromStock(ctx context.Context, materialID entity.ID, quantity float64) error {
	material, err := s.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	if err := material.RemoveFromStock(quantity); err != nil {
		return fmt.Errorf("ошибка списания со склада: %w", err)
	}

	if err := s.materialRepo.Update(ctx, material); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}

	return nil
}

// GetLowStockMaterials возвращает материалы с низким остатком
func (s *MaterialService) GetLowStockMaterials(ctx context.Context) ([]*entity.Material, error) {
	materials, err := s.materialRepo.GetLowStockMaterials(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения материалов с низким остатком: %w", err)
	}

	return materials, nil
}

// SearchMaterials ищет материалы по запросу
func (s *MaterialService) SearchMaterials(ctx context.Context, query string, limit, offset int) ([]*entity.Material, error) {
	materials, err := s.materialRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска материалов: %w", err)
	}

	return materials, nil
}

// SetMinStockQuantity устанавливает минимальный остаток для материала
func (s *MaterialService) SetMinStockQuantity(ctx context.Context, materialID entity.ID, minQuantity float64) error {
	material, err := s.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	if err := material.SetMinStockQuantity(minQuantity); err != nil {
		return fmt.Errorf("ошибка установки минимального остатка: %w", err)
	}

	if err := s.materialRepo.Update(ctx, material); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}

	return nil
} 