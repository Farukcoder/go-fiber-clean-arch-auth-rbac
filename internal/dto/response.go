package dto

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// PaginationResponse carries pagination metadata for list endpoints.
type PaginationResponse struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

// SuccessResponse builds a successful APIResponse.
func SuccessResponse(statusCode int, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:     true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

// ErrorResponse builds a failure APIResponse.
func ErrorResponse(statusCode int, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:     false,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
