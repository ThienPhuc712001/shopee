package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuery(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		q := NewQuery(0, 0)
		assert.Equal(t, 1, q.Page)
		assert.Equal(t, 20, q.Limit)
	})

	t.Run("valid values", func(t *testing.T) {
		q := NewQuery(5, 50)
		assert.Equal(t, 5, q.Page)
		assert.Equal(t, 50, q.Limit)
	})

	t.Run("limit exceeds max", func(t *testing.T) {
		// NewQuery caps limit at 100, defaults to 20 if exceeded
		q := NewQuery(1, 200)
		assert.Equal(t, 1, q.Page)
		assert.Equal(t, 20, q.Limit) // Limit > 100 defaults to 20
	})
}

func TestQuery_Default(t *testing.T) {
	t.Run("applies defaults", func(t *testing.T) {
		q := &Query{Page: 0, Limit: 0, SortOrder: ""}
		q.Default()
		assert.Equal(t, 1, q.Page)
		assert.Equal(t, 20, q.Limit)
		assert.Equal(t, "desc", q.SortOrder)
	})

	t.Run("preserves valid values", func(t *testing.T) {
		q := &Query{Page: 3, Limit: 50, SortOrder: "asc"}
		q.Default()
		assert.Equal(t, 3, q.Page)
		assert.Equal(t, 50, q.Limit)
		assert.Equal(t, "asc", q.SortOrder)
	})
}

func TestQuery_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"page 1, limit 20", 1, 20, 0},
		{"page 2, limit 20", 2, 20, 20},
		{"page 3, limit 20", 3, 20, 40},
		{"page 1, limit 50", 1, 50, 0},
		{"page 5, limit 50", 5, 50, 200},
		{"page 10, limit 100", 10, 100, 900},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{Page: tt.page, Limit: tt.limit}
			assert.Equal(t, tt.expected, q.Offset())
		})
	}
}

func TestNewResult(t *testing.T) {
	t.Run("first page", func(t *testing.T) {
		r := NewResult(1, 20, 300)
		assert.Equal(t, 1, r.Page)
		assert.Equal(t, 20, r.Limit)
		assert.Equal(t, int64(300), r.Total)
		assert.Equal(t, 15, r.TotalPages)
		assert.False(t, r.HasPrev)
		assert.True(t, r.HasNext)
		assert.Nil(t, r.PrevPage)
		assert.NotNil(t, r.NextPage)
		assert.Equal(t, 2, *r.NextPage)
	})

	t.Run("middle page", func(t *testing.T) {
		r := NewResult(5, 20, 300)
		assert.Equal(t, 5, r.Page)
		assert.Equal(t, 20, r.Limit)
		assert.Equal(t, int64(300), r.Total)
		assert.Equal(t, 15, r.TotalPages)
		assert.True(t, r.HasPrev)
		assert.True(t, r.HasNext)
		assert.NotNil(t, r.PrevPage)
		assert.NotNil(t, r.NextPage)
		assert.Equal(t, 4, *r.PrevPage)
		assert.Equal(t, 6, *r.NextPage)
	})

	t.Run("last page", func(t *testing.T) {
		r := NewResult(15, 20, 300)
		assert.Equal(t, 15, r.Page)
		assert.Equal(t, 20, r.Limit)
		assert.Equal(t, int64(300), r.Total)
		assert.Equal(t, 15, r.TotalPages)
		assert.True(t, r.HasPrev)
		assert.False(t, r.HasNext)
		assert.NotNil(t, r.PrevPage)
		assert.Nil(t, r.NextPage)
		assert.Equal(t, 14, *r.PrevPage)
	})

	t.Run("single page", func(t *testing.T) {
		r := NewResult(1, 20, 15)
		assert.Equal(t, 1, r.Page)
		assert.Equal(t, 20, r.Limit)
		assert.Equal(t, int64(15), r.Total)
		assert.Equal(t, 1, r.TotalPages)
		assert.False(t, r.HasPrev)
		assert.False(t, r.HasNext)
		assert.Nil(t, r.PrevPage)
		assert.Nil(t, r.NextPage)
	})

	t.Run("empty result", func(t *testing.T) {
		r := NewResult(1, 20, 0)
		assert.Equal(t, 1, r.Page)
		assert.Equal(t, 20, r.Limit)
		assert.Equal(t, int64(0), r.Total)
		assert.Equal(t, 0, r.TotalPages)
		assert.False(t, r.HasPrev)
		assert.False(t, r.HasNext)
	})

	t.Run("exact multiple", func(t *testing.T) {
		r := NewResult(1, 20, 100)
		assert.Equal(t, 5, r.TotalPages)
	})

	t.Run("non-exact multiple", func(t *testing.T) {
		r := NewResult(1, 20, 101)
		assert.Equal(t, 6, r.TotalPages)
	})
}

func TestNewResultFromQuery(t *testing.T) {
	q := NewQuery(3, 50)
	r := NewResultFromQuery(q, 500)

	assert.Equal(t, 3, r.Page)
	assert.Equal(t, 50, r.Limit)
	assert.Equal(t, int64(500), r.Total)
	assert.Equal(t, 10, r.TotalPages)
}

func TestValidatePageLimit(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		limit     int
		maxLimit  int
		wantPage  int
		wantLimit int
	}{
		{"valid values", 5, 50, 100, 5, 50},
		{"page too low", 0, 20, 100, 1, 20},
		{"limit too low", 1, 0, 100, 1, 20},
		{"limit too high", 1, 200, 100, 1, 100},
		{"negative page", -5, 20, 100, 1, 20},
		{"negative limit", 1, -10, 100, 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, limit := ValidatePageLimit(tt.page, tt.limit, tt.maxLimit)
			assert.Equal(t, tt.wantPage, page)
			assert.Equal(t, tt.wantLimit, limit)
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		limit    int
		expected int
	}{
		{"zero total", 0, 20, 0},
		{"exact multiple", 100, 20, 5},
		{"non-exact multiple", 101, 20, 6},
		{"single item", 1, 20, 1},
		{"limit zero", 100, 0, 5}, // defaults to 20
		{"large total", 10000, 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTotalPages(tt.total, tt.limit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaginate(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"page 1, limit 20", 1, 20, 0},
		{"page 2, limit 20", 2, 20, 20},
		{"page 5, limit 50", 5, 50, 200},
		{"page 10, limit 100", 10, 100, 900},
		{"invalid page", 0, 20, 0},    // defaults to page 1
		{"invalid limit", 1, 0, 0},    // defaults to limit 20
		{"limit too high", 1, 200, 0}, // defaults to limit 100
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := Paginate(tt.page, tt.limit)
			assert.Equal(t, tt.expected, offset)
		})
	}
}

func TestPagination_Integration(t *testing.T) {
	// Test complete pagination flow
	q := NewQuery(3, 25)
	q.Default()

	offset := q.Offset()
	assert.Equal(t, 50, offset)

	result := NewResultFromQuery(q, 234)

	assert.Equal(t, 3, result.Page)
	assert.Equal(t, 25, result.Limit)
	assert.Equal(t, int64(234), result.Total)
	assert.Equal(t, 10, result.TotalPages) // ceil(234/25) = 10
	assert.True(t, result.HasPrev)
	assert.True(t, result.HasNext) // Page 3 of 10, should have next

	// Verify next page
	assert.NotNil(t, result.NextPage)
	assert.Equal(t, 4, *result.NextPage)
}
