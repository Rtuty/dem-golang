package usecases

import (
	"testing"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductUseCaseTestSuite struct {
	suite.Suite
	productRepo  *mocks.MockProductRepository
	materialRepo *mocks.MockMaterialRepository
	useCase      *ProductUseCase
}

func (suite *ProductUseCaseTestSuite) SetupTest() {
	suite.productRepo = new(mocks.MockProductRepository)
	suite.materialRepo = new(mocks.MockMaterialRepository)
	suite.useCase = NewProductUseCase(suite.productRepo, suite.materialRepo)
}

func (suite *ProductUseCaseTestSuite) TestGetAllProducts() {
	// Подготовка данных
	products := []entities.Product{
		{
			ID:      1,
			Article: "ART001",
			Name:    "Обои винил",
			ProductType: &entities.ProductType{
				ID:          1,
				Name:        "Винил",
				Coefficient: 1.5,
			},
			Materials: []entities.ProductMaterial{
				{
					ID:              1,
					QuantityPerUnit: 2.0,
					Material: &entities.Material{
						ID:          1,
						CostPerUnit: 100.0,
					},
				},
			},
		},
	}

	// Настройка моков
	suite.productRepo.On("GetAll").Return(products, nil)

	// Выполнение
	result, err := suite.useCase.GetAllProducts()

	// Проверки
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), products[0].ID, result[0].ID)

	// Проверяем, что цена была рассчитана
	expectedPrice := 2.0 * 100.0 * 1.5 * 1.2 // = 360.0
	assert.NotNil(suite.T(), result[0].CalculatedPrice)
	assert.Equal(suite.T(), expectedPrice, *result[0].CalculatedPrice)

	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductUseCaseTestSuite) TestGetProductByID_Success() {
	// Подготовка данных
	productID := 1
	product := &entities.Product{
		ID:      productID,
		Article: "ART001",
		Name:    "Обои винил",
		ProductType: &entities.ProductType{
			ID:          1,
			Name:        "Винил",
			Coefficient: 1.5,
		},
		Materials: []entities.ProductMaterial{
			{
				ID:              1,
				QuantityPerUnit: 2.0,
				Material: &entities.Material{
					ID:          1,
					CostPerUnit: 100.0,
				},
			},
		},
	}

	// Настройка моков
	suite.productRepo.On("GetByID", productID).Return(product, nil)

	// Выполнение
	result, err := suite.useCase.GetProductByID(productID)

	// Проверки
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), product.ID, result.ID)

	// Проверяем, что цена была рассчитана
	assert.NotNil(suite.T(), result.CalculatedPrice)

	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductUseCaseTestSuite) TestGetProductByID_NotFound() {
	// Подготовка данных
	productID := 999

	// Настройка моков
	suite.productRepo.On("GetByID", productID).Return(nil, entities.NewNotFoundError("product", "999"))

	// Выполнение
	result, err := suite.useCase.GetProductByID(productID)

	// Проверки
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "продукции")

	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductUseCaseTestSuite) TestCreateProduct_Success() {
	// Подготовка данных
	product := &entities.Product{
		Article:         "ART002",
		Name:            "Новые обои",
		ProductTypeID:   1,
		MinPartnerPrice: 150.0,
	}

	productType := &entities.ProductType{
		ID:          1,
		Name:        "Винил",
		Coefficient: 1.5,
	}

	// Настройка моков
	suite.productRepo.On("Create", product).Return(nil)
	suite.productRepo.On("GetProductTypeByID", 1).Return(productType, nil)

	// Выполнение
	err := suite.useCase.CreateProduct(product)

	// Проверки
	assert.NoError(suite.T(), err)

	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductUseCaseTestSuite) TestCreateProduct_ValidationError() {
	// Подготовка данных (невалидный продукт)
	product := &entities.Product{
		Article:         "", // Пустой артикул
		Name:            "Новые обои",
		MinPartnerPrice: 150.0,
	}

	// Выполнение
	err := suite.useCase.CreateProduct(product)

	// Проверки
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "валидации")

	// Мок не должен вызываться при ошибке валидации
	suite.productRepo.AssertNotCalled(suite.T(), "Create")
}

func (suite *ProductUseCaseTestSuite) TestUpdateProduct_Success() {
	// Подготовка данных - изменяем тип продукции чтобы вызвался GetProductTypeByID
	product := &entities.Product{
		ID:              1,
		Article:         "ART001",
		Name:            "Обновленные обои",
		ProductTypeID:   2, // Меняем тип с 1 на 2
		MinPartnerPrice: 200.0,
	}

	existingProduct := &entities.Product{
		ID:              1,
		Article:         "ART001",
		Name:            "Старые обои",
		ProductTypeID:   1, // Старый тип
		MinPartnerPrice: 150.0,
	}

	productType := &entities.ProductType{
		ID:          2,
		Name:        "Бумага",
		Coefficient: 1.2,
	}

	// Настройка моков
	suite.productRepo.On("GetByID", 1).Return(existingProduct, nil)
	suite.productRepo.On("GetProductTypeByID", 2).Return(productType, nil)
	suite.productRepo.On("Update", product).Return(nil)

	// Выполнение
	err := suite.useCase.UpdateProduct(product)

	// Проверки
	assert.NoError(suite.T(), err)

	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductUseCaseTestSuite) TestDeleteProduct_Success() {
	// Подготовка данных
	productID := 1

	existingProduct := &entities.Product{
		ID:              1,
		Article:         "ART001",
		Name:            "Удаляемые обои",
		ProductTypeID:   1,
		MinPartnerPrice: 150.0,
	}

	// Настройка моков
	suite.productRepo.On("GetByID", productID).Return(existingProduct, nil)
	suite.productRepo.On("Delete", productID).Return(nil)

	// Выполнение
	err := suite.useCase.DeleteProduct(productID)

	// Проверки
	assert.NoError(suite.T(), err)

	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductUseCaseTestSuite) TestGetProductTypes() {
	// Подготовка данных
	productTypes := []entities.ProductType{
		{
			ID:          1,
			Name:        "Винил",
			Coefficient: 1.5,
		},
		{
			ID:          2,
			Name:        "Бумага",
			Coefficient: 1.2,
		},
	}

	// Настройка моков
	suite.productRepo.On("GetProductTypes").Return(productTypes, nil)

	// Выполнение
	result, err := suite.useCase.GetProductTypes()

	// Проверки
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), productTypes[0].Name, result[0].Name)
	assert.Equal(suite.T(), productTypes[1].Name, result[1].Name)

	suite.productRepo.AssertExpectations(suite.T())
}

func TestProductUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUseCaseTestSuite))
}

// Дополнительные тесты для edge cases
func TestProductUseCase_CalculatePriceEdgeCases(t *testing.T) {
	productRepo := new(mocks.MockProductRepository)
	materialRepo := new(mocks.MockMaterialRepository)
	useCase := NewProductUseCase(productRepo, materialRepo)

	t.Run("Продукт без типа", func(t *testing.T) {
		product := &entities.Product{
			ID:            1,
			ProductTypeID: 1,
			ProductType:   nil,
		}

		productType := &entities.ProductType{
			ID:          1,
			Name:        "Винил",
			Coefficient: 1.5,
		}

		// Настройка моков
		productRepo.On("GetProductTypeByID", 1).Return(productType, nil)
		productRepo.On("GetMaterialsForProduct", 1).Return([]entities.ProductMaterial{}, nil)

		price, err := useCase.calculateProductPrice(product)
		assert.NoError(t, err) // Теперь должно работать корректно
		assert.Equal(t, 0.0, price)

		productRepo.AssertExpectations(t)
	})

	t.Run("Продукт без материалов", func(t *testing.T) {
		product := &entities.Product{
			ID: 1,
			ProductType: &entities.ProductType{
				Coefficient: 1.5,
			},
			Materials: []entities.ProductMaterial{},
		}

		price, err := useCase.calculateProductPrice(product)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, price)
	})

	t.Run("Материал без цены", func(t *testing.T) {
		product := &entities.Product{
			ID: 1,
			ProductType: &entities.ProductType{
				Coefficient: 1.5,
			},
			Materials: []entities.ProductMaterial{
				{
					QuantityPerUnit: 1.0,
					Material: &entities.Material{
						CostPerUnit: 0,
					},
				},
			},
		}

		price, err := useCase.calculateProductPrice(product)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, price)
	})
}
