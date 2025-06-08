package tests

import (
	"testing"

	"wallpaper-system/internal/domain/entities"
	"wallpaper-system/internal/domain/mocks"
	"wallpaper-system/internal/usecases"
)

// BenchmarkProductCalculatePrice тестирует производительность расчета цены продукции
func BenchmarkProductCalculatePrice(b *testing.B) {
	product := &entities.Product{
		ID:   1,
		Name: "Обои тест",
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
					Name:        "Материал 1",
					CostPerUnit: 100.0,
				},
			},
			{
				ID:              2,
				QuantityPerUnit: 1.5,
				Material: &entities.Material{
					ID:          2,
					Name:        "Материал 2",
					CostPerUnit: 50.0,
				},
			},
			{
				ID:              3,
				QuantityPerUnit: 0.5,
				Material: &entities.Material{
					ID:          3,
					Name:        "Материал 3",
					CostPerUnit: 200.0,
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		product.CalculatePrice()
	}
}

// BenchmarkMaterialCalculateRequiredQuantity тестирует производительность расчета материала
func BenchmarkMaterialCalculateRequiredQuantity(b *testing.B) {
	material := &entities.Material{
		ID:          1,
		Name:        "Тестовый материал",
		CostPerUnit: 100.0,
	}

	baseQuantity := 10.5
	wastePercentage := 15.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		material.CalculateRequiredQuantity(baseQuantity, wastePercentage)
	}
}

// BenchmarkProductValidation тестирует производительность валидации продукции
func BenchmarkProductValidation(b *testing.B) {
	product := &entities.Product{
		Article:         "ART001",
		Name:            "Тестовые обои",
		MinPartnerPrice: 150.0,
		RollWidth:       func() *float64 { f := 1.06; return &f }(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		product.Validate()
	}
}

// BenchmarkMaterialValidation тестирует производительность валидации материала
func BenchmarkMaterialValidation(b *testing.B) {
	material := &entities.Material{
		Article:         "MAT001",
		Name:            "Тестовый материал",
		CostPerUnit:     100.0,
		PackageQuantity: 10.0,
		StockQuantity:   50.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		material.Validate()
	}
}

// BenchmarkCalculatorUseCase тестирует производительность калькулятора
func BenchmarkCalculatorUseCase(b *testing.B) {
	materialRepo := new(mocks.MockMaterialRepository)
	useCase := usecases.NewCalculatorUseCase(materialRepo)

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
	materialRepo.On("GetProductTypeByID", 1).Return(productType, nil)
	materialRepo.On("GetMaterialTypeByID", 1).Return(materialType, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useCase.CalculateRequiredMaterial(request)
	}
}

// BenchmarkProductUseCaseGetAllProducts тестирует производительность получения всех продуктов
func BenchmarkProductUseCaseGetAllProducts(b *testing.B) {
	productRepo := new(mocks.MockProductRepository)
	materialRepo := new(mocks.MockMaterialRepository)
	useCase := usecases.NewProductUseCase(productRepo, materialRepo)

	// Создаем множество тестовых продуктов
	products := make([]entities.Product, 1000)
	for i := 0; i < 1000; i++ {
		products[i] = entities.Product{
			ID:      i + 1,
			Article: "ART" + string(rune(i)),
			Name:    "Продукт " + string(rune(i)),
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
	}

	productRepo.On("GetAll").Return(products, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useCase.GetAllProducts()
	}
}

// BenchmarkValidationRequest тестирует производительность валидации запросов
func BenchmarkValidationRequest(b *testing.B) {
	request := &entities.MaterialCalculationRequest{
		ProductTypeID:   1,
		MaterialTypeID:  1,
		ProductQuantity: 10,
		ProductParam1:   2.0,
		ProductParam2:   1.5,
		MaterialInStock: 5.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request.Validate()
	}
}

// BenchmarkParallelProductCalculatePrice тестирует производительность в параллельном режиме
func BenchmarkParallelProductCalculatePrice(b *testing.B) {
	product := &entities.Product{
		ID:   1,
		Name: "Обои тест",
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

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			product.CalculatePrice()
		}
	})
}

// BenchmarkMultipleValidations тестирует валидацию множественных объектов
func BenchmarkMultipleValidations(b *testing.B) {
	products := make([]*entities.Product, 100)
	for i := 0; i < 100; i++ {
		products[i] = &entities.Product{
			Article:         "ART" + string(rune(i)),
			Name:            "Продукт " + string(rune(i)),
			MinPartnerPrice: 100.0 + float64(i),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, product := range products {
			product.Validate()
		}
	}
}
