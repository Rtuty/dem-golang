package mocks

import (
	"wallpaper-system/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

// MockCalculatorUseCase - мок для CalculatorUseCase
type MockCalculatorUseCase struct {
	mock.Mock
}

// CalculateRequiredMaterial рассчитывает необходимое количество материала
func (m *MockCalculatorUseCase) CalculateRequiredMaterial(request *entities.MaterialCalculationRequest) (int, error) {
	args := m.Called(request)
	return args.Int(0), args.Error(1)
}

// MockMaterialUseCase - мок для MaterialUseCase
type MockMaterialUseCase struct {
	mock.Mock
}

// GetAllMaterials возвращает список всех материалов
func (m *MockMaterialUseCase) GetAllMaterials() ([]entities.Material, error) {
	args := m.Called()
	return args.Get(0).([]entities.Material), args.Error(1)
}

// GetMaterialByID возвращает материал по ID
func (m *MockMaterialUseCase) GetMaterialByID(id int) (*entities.Material, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Material), args.Error(1)
}

// GetMaterialTypes возвращает все типы материалов
func (m *MockMaterialUseCase) GetMaterialTypes() ([]entities.MaterialType, error) {
	args := m.Called()
	return args.Get(0).([]entities.MaterialType), args.Error(1)
}

// GetMaterialsForProduct возвращает материалы для конкретной продукции
func (m *MockMaterialUseCase) GetMaterialsForProduct(productID int) ([]entities.Material, error) {
	args := m.Called(productID)
	return args.Get(0).([]entities.Material), args.Error(1)
}
