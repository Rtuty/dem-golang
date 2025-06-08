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

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (m *MockMaterialRepository) GetMaterialsForProduct(productID int) ([]entities.Material, error) {
	args := m.Called(productID)
	return args.Get(0).([]entities.Material), args.Error(1)
}
