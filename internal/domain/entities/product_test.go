package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_CalculatePrice(t *testing.T) {
	tests := []struct {
		name     string
		product  *Product
		expected float64
	}{
		{
			name: "Расчет цены с материалами и коэффициентом",
			product: &Product{
				ID:   1,
				Name: "Тестовые обои",
				ProductType: &ProductType{
					ID:          1,
					Name:        "Винил",
					Coefficient: 1.5,
				},
				Materials: []ProductMaterial{
					{
						ID:              1,
						QuantityPerUnit: 2.0,
						Material: &Material{
							ID:          1,
							Name:        "Винил основа",
							CostPerUnit: 100.0,
						},
					},
					{
						ID:              2,
						QuantityPerUnit: 0.5,
						Material: &Material{
							ID:          2,
							Name:        "Краска",
							CostPerUnit: 200.0,
						},
					},
				},
			},
			expected: 540.0, // (2*100 + 0.5*200) * 1.5 * 1.2 = 300 * 1.5 * 1.2 = 540
		},
		{
			name: "Продукция без материалов",
			product: &Product{
				ID:   2,
				Name: "Пустые обои",
				ProductType: &ProductType{
					ID:          1,
					Coefficient: 1.2,
				},
				Materials: []ProductMaterial{},
			},
			expected: 0,
		},
		{
			name: "Продукция без типа",
			product: &Product{
				ID:          3,
				Name:        "Обои без типа",
				ProductType: nil,
				Materials: []ProductMaterial{
					{
						ID:              1,
						QuantityPerUnit: 1.0,
						Material: &Material{
							ID:          1,
							CostPerUnit: 100.0,
						},
					},
				},
			},
			expected: 0,
		},
		{
			name: "Материал без цены",
			product: &Product{
				ID:   4,
				Name: "Обои с бесплатным материалом",
				ProductType: &ProductType{
					ID:          1,
					Coefficient: 1.0,
				},
				Materials: []ProductMaterial{
					{
						ID:              1,
						QuantityPerUnit: 1.0,
						Material: &Material{
							ID:          1,
							CostPerUnit: 0,
						},
					},
				},
			},
			expected: 0, // 0 * 1.0 * 1.2 = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.CalculatePrice()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product *Product
		wantErr bool
		errType string
	}{
		{
			name: "Валидная продукция",
			product: &Product{
				Article:         "ART001",
				Name:            "Тестовые обои",
				MinPartnerPrice: 100.0,
				RollWidth:       func() *float64 { f := 1.06; return &f }(),
			},
			wantErr: false,
		},
		{
			name: "Пустой артикул",
			product: &Product{
				Article:         "",
				Name:            "Тестовые обои",
				MinPartnerPrice: 100.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Пустое название",
			product: &Product{
				Article:         "ART001",
				Name:            "",
				MinPartnerPrice: 100.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Отрицательная цена",
			product: &Product{
				Article:         "ART001",
				Name:            "Тестовые обои",
				MinPartnerPrice: -100.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Отрицательная ширина рулона",
			product: &Product{
				Article:         "ART001",
				Name:            "Тестовые обои",
				MinPartnerPrice: 100.0,
				RollWidth:       func() *float64 { f := -1.0; return &f }(),
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Нулевая ширина рулона допустима",
			product: &Product{
				Article:         "ART001",
				Name:            "Тестовые обои",
				MinPartnerPrice: 100.0,
				RollWidth:       func() *float64 { f := 0.0; return &f }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != "" {
					validationErr, ok := err.(*ValidationError)
					assert.True(t, ok, "Ожидалась ValidationError")
					assert.Equal(t, tt.errType, validationErr.DomainError.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
