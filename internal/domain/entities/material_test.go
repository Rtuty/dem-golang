package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaterial_CalculateRequiredQuantity(t *testing.T) {
	tests := []struct {
		name            string
		material        *Material
		baseQuantity    float64
		wastePercentage float64
		expected        int
		expectError     bool
	}{
		{
			name: "Нормальный расчет с отходами",
			material: &Material{
				ID:          1,
				Name:        "Винил",
				CostPerUnit: 100.0,
			},
			baseQuantity:    10.0,
			wastePercentage: 15.0,
			expected:        12, // ceil(10 * (1 + 15/100)) = ceil(11.5) = 12
			expectError:     false,
		},
		{
			name: "Расчет без отходов",
			material: &Material{
				ID:          2,
				Name:        "Клей",
				CostPerUnit: 50.0,
			},
			baseQuantity:    5.0,
			wastePercentage: 0.0,
			expected:        5, // ceil(5 * (1 + 0/100)) = ceil(5) = 5
			expectError:     false,
		},
		{
			name: "Расчет с дробным базовым количеством",
			material: &Material{
				ID:          3,
				Name:        "Краска",
				CostPerUnit: 200.0,
			},
			baseQuantity:    2.3,
			wastePercentage: 10.0,
			expected:        3, // ceil(2.3 * (1 + 10/100)) = ceil(2.53) = 3
			expectError:     false,
		},
		{
			name: "Большой процент отходов",
			material: &Material{
				ID:          4,
				Name:        "Бумага",
				CostPerUnit: 30.0,
			},
			baseQuantity:    1.0,
			wastePercentage: 50.0,
			expected:        2, // ceil(1 * (1 + 50/100)) = ceil(1.5) = 2
			expectError:     false,
		},
		{
			name: "Нулевое базовое количество",
			material: &Material{
				ID:          5,
				Name:        "Тест",
				CostPerUnit: 100.0,
			},
			baseQuantity:    0.0,
			wastePercentage: 10.0,
			expected:        0, // ceil(0 * (1 + 10/100)) = ceil(0) = 0
			expectError:     false,
		},
		{
			name: "Отрицательное базовое количество",
			material: &Material{
				ID:          6,
				Name:        "Негатив",
				CostPerUnit: 100.0,
			},
			baseQuantity:    -5.0,
			wastePercentage: 10.0,
			expected:        0,
			expectError:     true,
		},
		{
			name: "Отрицательный процент отходов",
			material: &Material{
				ID:          7,
				Name:        "Тест2",
				CostPerUnit: 100.0,
			},
			baseQuantity:    5.0,
			wastePercentage: -10.0,
			expected:        0,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.material.CalculateRequiredQuantity(tt.baseQuantity, tt.wastePercentage)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMaterialCalculationRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *MaterialCalculationRequest
		wantErr bool
		errType string
	}{
		{
			name: "Валидный запрос",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: 10,
				ProductParam1:   5.0,
				ProductParam2:   3.0,
				MaterialInStock: 2.0,
			},
			wantErr: false,
		},
		{
			name: "Нулевой ID типа продукции",
			request: &MaterialCalculationRequest{
				ProductTypeID:   0,
				MaterialTypeID:  1,
				ProductQuantity: 10,
				ProductParam1:   5.0,
				ProductParam2:   3.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Нулевой ID типа материала",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  0,
				ProductQuantity: 10,
				ProductParam1:   5.0,
				ProductParam2:   3.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Нулевое количество продукции",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: 0,
				ProductParam1:   5.0,
				ProductParam2:   3.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Отрицательное количество продукции",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: -5,
				ProductParam1:   5.0,
				ProductParam2:   3.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Отрицательный первый параметр",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: 10,
				ProductParam1:   -1.0,
				ProductParam2:   3.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Отрицательный второй параметр",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: 10,
				ProductParam1:   5.0,
				ProductParam2:   -2.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Отрицательный остаток материала",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: 10,
				ProductParam1:   5.0,
				ProductParam2:   3.0,
				MaterialInStock: -1.0,
			},
			wantErr: true,
			errType: "VALIDATION_ERROR",
		},
		{
			name: "Нулевые параметры продукции допустимы",
			request: &MaterialCalculationRequest{
				ProductTypeID:   1,
				MaterialTypeID:  1,
				ProductQuantity: 1,
				ProductParam1:   0.0,
				ProductParam2:   0.0,
				MaterialInStock: 0.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

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
