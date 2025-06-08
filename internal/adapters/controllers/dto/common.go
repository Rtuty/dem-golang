package dto

// SuccessResponse представляет успешный ответ API
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// NewSuccessResponse создает новый успешный ответ
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse создает новый ответ с ошибкой
func NewErrorResponse(error string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   error,
	}
}
