package response

import "ecommerce/pkg/pagination"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response with full pagination metadata
type PaginatedResponse struct {
	Success    bool             `json:"success"`
	Data       interface{}      `json:"data"`
	Pagination *pagination.Result `json:"pagination"`
	Message    string           `json:"message,omitempty"`
}

// Meta contains pagination metadata (legacy format)
type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

// GenericPaginatedResponse is a generic paginated response for type safety
type GenericPaginatedResponse[T any] struct {
	Success    bool             `json:"success"`
	Data       []T              `json:"data"`
	Pagination *pagination.Result `json:"pagination"`
	Message    string           `json:"message,omitempty"`
}

// Success returns a success response
func Success(data interface{}, message string) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// SuccessWithMessage returns a success response with message only
func SuccessWithMessage(message string) Response {
	return Response{
		Success: true,
		Message: message,
	}
}

// Error returns an error response
func Error(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// BadRequest returns a bad request response
func BadRequest(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// Unauthorized returns an unauthorized response
func Unauthorized(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// Forbidden returns a forbidden response
func Forbidden(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// NotFound returns a not found response
func NotFound(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// InternalError returns an internal server error response
func InternalError(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// Paginated returns a paginated response with full pagination metadata
func Paginated(data interface{}, total int64, page, perPage int, message string) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination.NewResult(page, perPage, total),
		Message:    message,
	}
}

// PaginatedWithMeta returns a paginated response using legacy Meta format
func PaginatedWithMeta(data interface{}, total int64, page, perPage int, message string) PaginatedResponse {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: &pagination.Result{
			Page:       page,
			Limit:      perPage,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Message: message,
	}
}

// NewGenericPaginatedResponse creates a new generic paginated response
func NewGenericPaginatedResponse[T any](data []T, total int64, page, limit int, message string) *GenericPaginatedResponse[T] {
	return &GenericPaginatedResponse[T]{
		Success:    true,
		Data:       data,
		Pagination: pagination.NewResult(page, limit, total),
		Message:    message,
	}
}
