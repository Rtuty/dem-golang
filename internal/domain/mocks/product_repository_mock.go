package mocks

import (
	"wallpaper-system/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

// MockProductRepository - мок для интерфейса ProductRepository
type MockProductRepository struct {
	mock.Mock
}

// GetAll возвращает список всей продукции
func (m *MockProductRepository) GetAll() ([]entities.Product, error) {
	args := m.Called()
	return args.Get(0).([]entities.Product), args.Error(1)
}

// GetByID возвращает продукцию по ID
func (m *MockProductRepository) GetByID(id int) (*entities.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

// Create создает новую продукцию
func (m *MockProductRepository) Create(product *entities.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// Update обновляет продукцию
func (m *MockProductRepository) Update(product *entities.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// Delete удаляет продукцию
func (m *MockProductRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetProductTypes возвращает все типы продукции
func (m *MockProductRepository) GetProductTypes() ([]entities.ProductType, error) {
	args := m.Called()
	return args.Get(0).([]entities.ProductType), args.Error(1)
}

// GetProductTypeByID возвращает тип продукции по ID
func (m *MockProductRepository) GetProductTypeByID(id int) (*entities.ProductType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProductType), args.Error(1)
}

// GetMaterialsForProduct возвращает материалы для продукции
func (m *MockProductRepository) GetMaterialsForProduct(productID int) ([]entities.ProductMaterial, error) {
	args := m.Called(productID)
	return args.Get(0).([]entities.ProductMaterial), args.Error(1)
}
