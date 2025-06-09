package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"wallpaper-system/internal/adapters/controllers/dto"
	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/usecases/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductControllerTestSuite struct {
	suite.Suite
	productUseCase  *mocks.MockProductUseCase
	materialUseCase *mocks.MockMaterialUseCase
	controller      *ProductController
	router          *gin.Engine
}

func (suite *ProductControllerTestSuite) SetupTest() {
	suite.productUseCase = new(mocks.MockProductUseCase)
	suite.materialUseCase = new(mocks.MockMaterialUseCase)
	suite.controller = NewProductController(suite.productUseCase, suite.materialUseCase)

	// Настройка Gin в тестовом режиме
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Настройка маршрутов
	v1 := suite.router.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("/", suite.controller.GetProducts)
			products.GET("/:id", suite.controller.GetProductByID)
			products.POST("/", suite.controller.CreateProduct)
			products.PUT("/:id", suite.controller.UpdateProduct)
			products.DELETE("/:id", suite.controller.DeleteProduct)
		}
	}
}

func (suite *ProductControllerTestSuite) TestGetProducts_Success() {
	// Подготовка данных
	products := []entities.Product{
		{
			ID:      1,
			Article: "ART001",
			Name:    "Обои винил",
			ProductType: &entities.ProductType{
				ID:   1,
				Name: "Винил",
			},
		},
		{
			ID:      2,
			Article: "ART002",
			Name:    "Обои бумага",
			ProductType: &entities.ProductType{
				ID:   2,
				Name: "Бумага",
			},
		},
	}

	// Настройка мока
	suite.productUseCase.On("GetAllProducts").Return(products, nil)

	// Выполнение запроса
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)

	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetProducts_InternalError() {
	// Настройка мока для возврата ошибки
	suite.productUseCase.On("GetAllProducts").Return([]entities.Product{}, assert.AnError)

	// Выполнение запроса
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.Error)

	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetProduct_Success() {
	// Подготовка данных
	productID := 1
	product := &entities.Product{
		ID:      productID,
		Article: "ART001",
		Name:    "Обои винил",
		ProductType: &entities.ProductType{
			ID:   1,
			Name: "Винил",
		},
	}

	// Настройка мока
	suite.productUseCase.On("GetProductByID", productID).Return(product, nil)

	// Выполнение запроса
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+strconv.Itoa(productID), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)

	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetProduct_NotFound() {
	// Подготовка данных
	productID := 999

	// Настройка мока
	suite.productUseCase.On("GetProductByID", productID).Return(nil, entities.NewNotFoundError("product", "999"))

	// Выполнение запроса
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+strconv.Itoa(productID), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.Error)

	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetProduct_InvalidID() {
	// Выполнение запроса с невалидным ID
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/invalid", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Error, "Некорректный ID")

	// Мок не должен вызываться при невалидном ID
	suite.productUseCase.AssertNotCalled(suite.T(), "GetProductByID")
}

func (suite *ProductControllerTestSuite) TestCreateProduct_Success() {
	// Подготовка данных
	createRequest := dto.CreateProductRequest{
		Article:         "ART003",
		ProductTypeID:   1,
		Name:            "Новые обои",
		Description:     "Описание",
		MinPartnerPrice: 150.0,
		RollWidth:       float64Ptr(1.06),
	}

	// Настройка мока
	suite.productUseCase.On("CreateProduct", mock.AnythingOfType("*entities.Product")).Return(nil)

	// Подготовка запроса
	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestCreateProduct_ValidationError() {
	// Подготовка данных с ошибкой валидации (отсутствующий обязательный article)
	requestData := `{"product_type_id": 1, "name": "Новые обои", "min_partner_price": 150.0}`

	// Подготовка запроса (без article - обязательного поля)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", bytes.NewBufferString(requestData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Error, "Некорректные данные запроса")

	// Мок не должен вызываться при ошибке валидации на уровне контроллера
	suite.productUseCase.AssertNotCalled(suite.T(), "CreateProduct")
}

func (suite *ProductControllerTestSuite) TestCreateProduct_InvalidJSON() {
	// Выполнение запроса с невалидным JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Error, "Некорректные данные запроса")

	// Мок не должен вызываться при невалидном JSON
	suite.productUseCase.AssertNotCalled(suite.T(), "CreateProduct")
}

func (suite *ProductControllerTestSuite) TestDeleteProduct_Success() {
	// Подготовка данных
	productID := 1

	// Настройка мока
	suite.productUseCase.On("DeleteProduct", productID).Return(nil)

	// Выполнение запроса
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+strconv.Itoa(productID), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestDeleteProduct_NotFound() {
	// Подготовка данных
	productID := 999

	// Настройка мока
	suite.productUseCase.On("DeleteProduct", productID).Return(entities.NewNotFoundError("product", "999"))

	// Выполнение запроса
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+strconv.Itoa(productID), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Проверки
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.Error)

	suite.productUseCase.AssertExpectations(suite.T())
}

func TestProductControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductControllerTestSuite))
}

// Вспомогательные функции для создания указателей
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
