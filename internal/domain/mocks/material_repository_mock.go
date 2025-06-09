package mocks

import (
	"wallpaper-system/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

// MockMaterialRepository - мок для интерфейса MaterialRepository
type MockMaterialRepository struct {
	mock.Mock
}

// GetAll возвращает список всех материалов
func (m *MockMaterialRepository) GetAll() ([]entities.Material, error) {
	args := m.Called()
	return args.Get(0).([]entities.Material), args.Error(1)
}

// GetByID возвращает материал по ID
func (m *MockMaterialRepository) GetByID(id int) (*entities.Material, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Material), args.Error(1)
}

// Create создает новый материал
func (m *MockMaterialRepository) Create(material *entities.Material) error {
	args := m.Called(material)
	return args.Error(0)
}

// Update обновляет существующий материал
func (m *MockMaterialRepository) Update(material *entities.Material) error {
	args := m.Called(material)
	return args.Error(0)
}

// Delete удаляет материал по ID
func (m *MockMaterialRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetMaterialTypeByID возвращает тип материала по ID
func (m *MockMaterialRepository) GetMaterialTypeByID(id int) (*entities.MaterialType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.MaterialType), args.Error(1)
}

// GetProductTypeByID возвращает тип продукции по ID
func (m *MockMaterialRepository) GetProductTypeByID(id int) (*entities.ProductType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductType), args.Error(1)
}

// GetMaterialTypes возвращает все типы материалов
func (m *MockMaterialRepository) GetMaterialTypes() ([]entities.MaterialType, error) {
	args := m.Called()
	return args.Get(0).([]entities.MaterialType), args.Error(1)
}

// GetMeasurementUnits возвращает все единицы измерения
func (m *MockMaterialRepository) GetMeasurementUnits() ([]entities.MeasurementUnit, error) {
	args := m.Called()
	return args.Get(0).([]entities.MeasurementUnit), args.Error(1)
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (m *MockMaterialRepository) GetMaterialsForProduct(productID int) ([]entities.Material, error) {
	args := m.Called(productID)
	return args.Get(0).([]entities.Material), args.Error(1)
}
