package pagination

import (
	"math"

	"gorm.io/gorm"
)

// Query represents pagination query parameters from API requests
type Query struct {
	// Page is the page number (1-indexed, default: 1)
	Page int `form:"page" binding:"min=1"`

	// Limit is the number of items per page (default: 20, max: 100)
	Limit int `form:"limit" binding:"min=1,max=100"`

	// SortBy is the field to sort by (default: created_at)
	SortBy string `form:"sort_by"`

	// SortOrder is the sort direction: "asc" or "desc" (default: desc)
	SortOrder string `form:"sort_order"`
}

// Result represents pagination metadata in API responses
type Result struct {
	// Page is the current page number
	Page int `json:"page"`

	// Limit is the number of items per page
	Limit int `json:"limit"`

	// Total is the total number of items in the database
	Total int64 `json:"total"`

	// TotalPages is the total number of pages
	TotalPages int `json:"total_pages"`

	// HasNext indicates if there is a next page
	HasNext bool `json:"has_next"`

	// HasPrev indicates if there is a previous page
	HasPrev bool `json:"has_prev"`

	// NextPage is the next page number (nil if no next page)
	NextPage *int `json:"next_page,omitempty"`

	// PrevPage is the previous page number (nil if no previous page)
	PrevPage *int `json:"prev_page,omitempty"`
}

// NewQuery creates a new Query with default values
func NewQuery(page, limit int) *Query {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return &Query{
		Page:      page,
		Limit:     limit,
		SortOrder: "desc",
	}
}

// Default applies default values to the query
func (q *Query) Default() *Query {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	return q
}

// Offset calculates the database offset based on page and limit
// Formula: offset = (page - 1) × limit
func (q *Query) Offset() int {
	return (q.Page - 1) * q.Limit
}

// Apply applies pagination to a GORM query
func (q *Query) Apply(query *gorm.DB) *gorm.DB {
	return query.Offset(q.Offset()).Limit(q.Limit)
}

// ApplyWithSort applies pagination and sorting to a GORM query
func (q *Query) ApplyWithSort(query *gorm.DB, defaultSort string) *gorm.DB {
	sortBy := q.SortBy
	if sortBy == "" {
		sortBy = defaultSort
	}

	sortOrder := q.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return query.Order(sortBy + " " + sortOrder).Offset(q.Offset()).Limit(q.Limit)
}

// NewResult creates a new pagination Result
func NewResult(page, limit int, total int64) *Result {
	// Calculate total pages using ceiling division
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	result := &Result{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	if result.HasNext {
		next := page + 1
		result.NextPage = &next
	}
	if result.HasPrev {
		prev := page - 1
		result.PrevPage = &prev
	}

	return result
}

// NewResultFromQuery creates a new pagination Result from a Query
func NewResultFromQuery(q *Query, total int64) *Result {
	return NewResult(q.Page, q.Limit, total)
}

// ValidatePageLimit validates and normalizes page and limit values
func ValidatePageLimit(page, limit, maxLimit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return page, limit
}

// CalculateTotalPages calculates total pages from total items and limit
func CalculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		limit = 20
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

// Paginate calculates the offset for database queries
// This is a convenience function for simple pagination
func Paginate(page int, limit int) (offset int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return (page - 1) * limit
}
