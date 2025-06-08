package entities

import "fmt"

// DomainError представляет доменную ошибку
type DomainError struct {
	Code    string
	Message string
	Field   string
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	DomainError
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		DomainError: DomainError{
			Code:    "VALIDATION_ERROR",
			Message: message,
			Field:   field,
		},
	}
}

// BusinessError представляет бизнес-ошибку
type BusinessError struct {
	DomainError
}

// NewBusinessError создает новую бизнес-ошибку
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{
		DomainError: DomainError{
			Code:    code,
			Message: message,
		},
	}
}

// NotFoundError представляет ошибку "не найдено"
type NotFoundError struct {
	DomainError
}

// NewNotFoundError создает новую ошибку "не найдено"
func NewNotFoundError(entity, id string) *NotFoundError {
	return &NotFoundError{
		DomainError: DomainError{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("%s с ID %s не найден", entity, id),
		},
	}
}
