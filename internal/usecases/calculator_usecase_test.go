package usecases

import (
	"testing"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CalculatorUseCaseTestSuite struct {
	suite.Suite
	materialRepo *mocks.MockMaterialRepository
	useCase      *CalculatorUseCase
}

func (suite *CalculatorUseCaseTestSuite) SetupTest() {
	suite.materialRepo = new(mocks.MockMaterialRepository)
	suite.useCase = NewCalculatorUseCase(suite.materialRepo)
}

func (suite *CalculatorUseCaseTestSuite) TestCalculateRequiredMaterial_Success() {
	// Подготовка данных
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   1,
		MaterialTypeID:  1,
		ProductQuantity: 10,
		ProductParam1:   2.0,
		ProductParam2:   1.5,
		MaterialInStock: 5.0,
	}

	productType := &entities.ProductType{
		ID:          1,
		Name:        "Винил",
		Coefficient: 1.5,
	}

	materialType := &entities.MaterialType{
		ID:              1,
		Name:            "Основа",
		WastePercentage: 10.0,
	}

	// Настройка моков
	suite.materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
	suite.materialRepo.On("GetMaterialTypeByID", 1).Return(materialType, nil)

	// Выполнение
	result, err := suite.useCase.CalculateRequiredMaterial(request)

	// Проверки
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), result, 0)

	// Проверяем математику:
	// materialPerUnit = 2.0 * 1.5 * 1.5 = 4.5
	// totalMaterial = 4.5 * 10 = 45.0
	// materialWithWaste = 45.0 * (1 + 10/100) = 49.5
	// requiredMaterial = 49.5 - 5.0 = 44.5
	// rounded up = 45
	expectedResult := 45
	assert.Equal(suite.T(), expectedResult, result)

	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *CalculatorUseCaseTestSuite) TestCalculateRequiredMaterial_ValidationError() {
	// Подготовка данных (невалидный запрос)
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   0, // Невалидный ID
		MaterialTypeID:  1,
		ProductQuantity: 10,
		ProductParam1:   2.0,
		ProductParam2:   1.5,
	}

	// Выполнение
	result, err := suite.useCase.CalculateRequiredMaterial(request)

	// Проверки
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, result)
	assert.Contains(suite.T(), err.Error(), "валидации")

	// Моки не должны вызываться при ошибке валидации
	suite.materialRepo.AssertNotCalled(suite.T(), "GetProductTypeByID")
	suite.materialRepo.AssertNotCalled(suite.T(), "GetMaterialTypeByID")
}

func (suite *CalculatorUseCaseTestSuite) TestCalculateRequiredMaterial_ProductTypeNotFound() {
	// Подготовка данных
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   999, // Несуществующий ID
		MaterialTypeID:  1,
		ProductQuantity: 10,
		ProductParam1:   2.0,
		ProductParam2:   1.5,
	}

	// Настройка моков
	suite.materialRepo.On("GetProductTypeByID", 999).Return(nil, entities.NewNotFoundError("product_type", "999"))

	// Выполнение
	result, err := suite.useCase.CalculateRequiredMaterial(request)

	// Проверки
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, result)
	assert.Contains(suite.T(), err.Error(), "тип продукции")

	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *CalculatorUseCaseTestSuite) TestCalculateRequiredMaterial_MaterialTypeNotFound() {
	// Подготовка данных
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   1,
		MaterialTypeID:  999, // Несуществующий ID
		ProductQuantity: 10,
		ProductParam1:   2.0,
		ProductParam2:   1.5,
	}

	productType := &entities.ProductType{
		ID:          1,
		Name:        "Винил",
		Coefficient: 1.5,
	}

	// Настройка моков
	suite.materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
	suite.materialRepo.On("GetMaterialTypeByID", 999).Return(nil, entities.NewNotFoundError("material_type", "999"))

	// Выполнение
	result, err := suite.useCase.CalculateRequiredMaterial(request)

	// Проверки
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, result)
	assert.Contains(suite.T(), err.Error(), "тип материала")

	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *CalculatorUseCaseTestSuite) TestCalculateRequiredMaterial_NoAdditionalMaterialNeeded() {
	// Подготовка данных - достаточно материала на складе
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   1,
		MaterialTypeID:  1,
		ProductQuantity: 1,
		ProductParam1:   1.0,
		ProductParam2:   1.0,
		MaterialInStock: 100.0, // Много материала на складе
	}

	productType := &entities.ProductType{
		ID:          1,
		Name:        "Винил",
		Coefficient: 1.0,
	}

	materialType := &entities.MaterialType{
		ID:              1,
		Name:            "Основа",
		WastePercentage: 0.0,
	}

	// Настройка моков
	suite.materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
	suite.materialRepo.On("GetMaterialTypeByID", 1).Return(materialType, nil)

	// Выполнение
	result, err := suite.useCase.CalculateRequiredMaterial(request)

	// Проверки
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, result) // Дополнительный материал не нужен

	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *CalculatorUseCaseTestSuite) TestCalculateRequiredMaterial_WithHighWastePercentage() {
	// Подготовка данных с высоким процентом отходов
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   1,
		MaterialTypeID:  1,
		ProductQuantity: 5,
		ProductParam1:   2.0,
		ProductParam2:   2.0,
		MaterialInStock: 0.0,
	}

	productType := &entities.ProductType{
		ID:          1,
		Name:        "Винил",
		Coefficient: 1.0,
	}

	materialType := &entities.MaterialType{
		ID:              1,
		Name:            "Основа",
		WastePercentage: 50.0, // Высокий процент отходов
	}

	// Настройка моков
	suite.materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
	suite.materialRepo.On("GetMaterialTypeByID", 1).Return(materialType, nil)

	// Выполнение
	result, err := suite.useCase.CalculateRequiredMaterial(request)

	// Проверки
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), result, 0)

	// Проверяем математику:
	// materialPerUnit = 2.0 * 2.0 * 1.0 = 4.0
	// totalMaterial = 4.0 * 5 = 20.0
	// materialWithWaste = 20.0 * (1 + 50/100) = 30.0
	// requiredMaterial = 30.0 - 0.0 = 30.0
	// rounded up = 30
	expectedResult := 30
	assert.Equal(suite.T(), expectedResult, result)

	suite.materialRepo.AssertExpectations(suite.T())
}

func TestCalculatorUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CalculatorUseCaseTestSuite))
}

// Дополнительные unit tests для специфических случаев
func TestCalculatorUseCase_EdgeCases(t *testing.T) {
	materialRepo := new(mocks.MockMaterialRepository)
	useCase := NewCalculatorUseCase(materialRepo)

	t.Run("Минимальные значения", func(t *testing.T) {
		request := &entities.MaterialCalculationRequest{
			ProductTypeID:   1,
			MaterialTypeID:  1,
			ProductQuantity: 1,
			ProductParam1:   0.1,
			ProductParam2:   0.1,
			MaterialInStock: 0.0,
		}

		productType := &entities.ProductType{
			ID:          1,
			Coefficient: 0.1,
		}

		materialType := &entities.MaterialType{
			ID:              1,
			WastePercentage: 0.1,
		}

		materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
		materialRepo.On("GetMaterialTypeByID", 1).Return(materialType, nil)

		result, err := useCase.CalculateRequiredMaterial(request)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, result, 0)

		materialRepo.AssertExpectations(t)
	})

	t.Run("Дробные параметры", func(t *testing.T) {
		request := &entities.MaterialCalculationRequest{
			ProductTypeID:   1,
			MaterialTypeID:  1,
			ProductQuantity: 3,
			ProductParam1:   1.33,
			ProductParam2:   2.67,
			MaterialInStock: 1.5,
		}

		productType := &entities.ProductType{
			ID:          1,
			Coefficient: 1.25,
		}

		materialType := &entities.MaterialType{
			ID:              1,
			WastePercentage: 7.5,
		}

		materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
		materialRepo.On("GetMaterialTypeByID", 1).Return(materialType, nil)

		result, err := useCase.CalculateRequiredMaterial(request)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, result, 0)

		materialRepo.AssertExpectations(t)
	})
}
