package models

// APIResponse is the standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains pagination and additional info
type Meta struct {
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
}

// Helper functions
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func SuccessWithMeta(data interface{}, total int) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Total: total,
		},
	}
}

func ErrorResponse(code string, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}
