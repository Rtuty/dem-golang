package mocks

import (
	"wallpaper-system/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

// MockProductUseCase - мок для ProductUseCase
type MockProductUseCase struct {
	mock.Mock
}

// GetAllProducts возвращает список всей продукции с рассчитанными ценами
func (m *MockProductUseCase) GetAllProducts() ([]entities.Product, error) {
	args := m.Called()
	return args.Get(0).([]entities.Product), args.Error(1)
}

// GetProductByID возвращает продукцию по ID с материалами и рассчитанной ценой
func (m *MockProductUseCase) GetProductByID(id int) (*entities.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Product), args.Error(1)
}

// CreateProduct создает новую продукцию
func (m *MockProductUseCase) CreateProduct(product *entities.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// UpdateProduct обновляет продукцию
func (m *MockProductUseCase) UpdateProduct(product *entities.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// DeleteProduct удаляет продукцию
func (m *MockProductUseCase) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetProductTypes возвращает все типы продукции
func (m *MockProductUseCase) GetProductTypes() ([]entities.ProductType, error) {
	args := m.Called()
	return args.Get(0).([]entities.ProductType), args.Error(1)
}
