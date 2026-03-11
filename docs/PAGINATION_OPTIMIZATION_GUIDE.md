# Performance Optimization Guide for E-Commerce Platform

## Table of Contents

1. [Performance Problem](#part-1---performance-problem)
2. [Pagination Concept](#part-2---pagination-concept)
3. [Pagination Response Format](#part-3---pagination-response-format)
4. [Database Indexing](#part-4---database-indexing)
5. [Golang Pagination Struct](#part-5---golang-pagination-struct)
6. [Pagination Helper](#part-6---pagination-helper)
7. [Product List API](#part-7---product-list-api)
8. [Search Optimization](#part-8---search-optimization)
9. [Database Query Optimization](#part-9---database-query-optimization)
10. [Simple Caching](#part-10---simple-caching)
11. [Load Handling](#part-11---load-handling)
12. [Example Implementation](#part-12---example-implementation)

---

## PART 1 - PERFORMANCE PROBLEM

### Why Performance Optimization is Critical in E-Commerce

Performance optimization is essential for e-commerce systems due to several key factors:

### 1. Large Product Catalogs

Modern e-commerce platforms often manage:
- **Thousands to millions of products** across multiple categories
- **Multiple product variants** (size, color, etc.) per product
- **Rich media content** (images, videos, 3D models)
- **Complex relationships** (categories, tags, attributes, reviews)

**Impact without optimization:**
```
SELECT * FROM products;  -- Returns 500,000+ rows
-- Memory usage: 2-5 GB
-- Query time: 30+ seconds
-- Network transfer: 100+ MB
```

### 2. High User Traffic

E-commerce platforms experience:
- **Concurrent users**: Hundreds to thousands simultaneously browsing
- **Peak traffic**: 10x normal load during sales/events
- **API requests per user**: 20-50 requests per session

**Example scenario:**
```
1,000 concurrent users × 10 product list requests each = 10,000 queries
Without pagination: 10,000 × 500,000 rows = 5 billion rows processed
With pagination (20 items): 10,000 × 20 rows = 200,000 rows processed
Improvement: 25,000x reduction in data processing
```

### 3. Database Load

Inefficient queries cause:
- **Connection pool exhaustion**: All database connections occupied
- **Lock contention**: Queries waiting for table locks
- **I/O bottleneck**: Disk reads overwhelming the system
- **Memory pressure**: Large result sets consuming RAM

### How Inefficient Queries Affect Performance

| Problem | Impact | Solution |
|---------|--------|----------|
| SELECT * queries | High memory, network usage | Select specific columns |
| No pagination | Large result sets | Implement LIMIT/OFFSET |
| Missing indexes | Full table scans | Add appropriate indexes |
| N+1 queries | Multiple round trips | Use eager loading (JOINs) |
| No caching | Repeated DB hits | Implement caching layer |

---

## PART 2 - PAGINATION CONCEPT

### What is Pagination?

Pagination is a technique to divide large datasets into smaller, manageable chunks (pages) that are fetched and displayed one at a time.

### API Request Format

```http
GET /api/products?page=1&limit=20
```

### Key Concepts

| Term | Definition | Example |
|------|------------|---------|
| **page** | Current page number (1-indexed) | `page=1` (first page) |
| **limit** | Number of items per page | `limit=20` (20 items/page) |
| **offset** | Number of records to skip | `offset=0` (start from beginning) |

### Offset Calculation Formula

```
offset = (page - 1) × limit
```

### Examples

```
Page 1, Limit 20: offset = (1 - 1) × 20 = 0
Page 2, Limit 20: offset = (2 - 1) × 20 = 20
Page 3, Limit 20: offset = (3 - 1) × 20 = 40
Page 10, Limit 50: offset = (10 - 1) × 50 = 450
```

### SQL Implementation

```sql
-- Page 1, 20 items per page
SELECT * FROM products 
WHERE status = 'active'
ORDER BY created_at DESC
OFFSET 0 ROWS
FETCH NEXT 20 ROWS ONLY;

-- Page 2, 20 items per page
SELECT * FROM products 
WHERE status = 'active'
ORDER BY created_at DESC
OFFSET 20 ROWS
FETCH NEXT 20 ROWS ONLY;
```

---

## PART 3 - PAGINATION RESPONSE FORMAT

### Standard API Response Structure

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Product 1",
        "price": 99.99,
        "image_url": "/images/product1.jpg"
      },
      {
        "id": 2,
        "name": "Product 2",
        "price": 149.99,
        "image_url": "/images/product2.jpg"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 300,
      "total_pages": 15,
      "has_next": true,
      "has_prev": false,
      "next_page": 2,
      "prev_page": null
    }
  }
}
```

### Field Explanations

| Field | Type | Description |
|-------|------|-------------|
| `page` | int | Current page number |
| `limit` | int | Items per page |
| `total` | int64 | Total number of items in database |
| `total_pages` | int | Total number of pages (ceil(total/limit)) |
| `has_next` | bool | Whether next page exists |
| `has_prev` | bool | Whether previous page exists |
| `next_page` | *int | Next page number (null if none) |
| `prev_page` | *int | Previous page number (null if none) |

### Compact Response Format (Alternative)

```json
{
  "success": true,
  "data": [],
  "page": 1,
  "limit": 20,
  "total": 300,
  "total_pages": 15
}
```

---

## PART 4 - DATABASE INDEXING

### Why Indexes are Important

Database indexes are data structures that improve the speed of data retrieval operations:

| Benefit | Description |
|---------|-------------|
| **Faster queries** | O(log n) vs O(n) lookup time |
| **Reduced I/O** | Fewer disk reads required |
| **Better sorting** | Pre-sorted index data |
| **Efficient joins** | Faster relationship lookups |

### Recommended Indexes for E-Commerce

```sql
-- ============================================
-- PRODUCT INDEXES
-- ============================================

-- Index on products.name for search functionality
CREATE INDEX IX_products_name ON products(name);

-- Composite index for category filtering + status
CREATE INDEX IX_products_category_id_status ON products(category_id, status);

-- Index for shop filtering
CREATE INDEX IX_products_shop_id ON products(shop_id);

-- Index for status filtering (active products)
CREATE INDEX IX_products_status ON products(status);

-- Index for featured products
CREATE INDEX IX_products_is_featured ON products(is_featured) WHERE status = 'active';

-- Index for best sellers (sorted by sold_count)
CREATE INDEX IX_products_sold_count ON products(sold_count DESC);

-- Index for price range queries
CREATE INDEX IX_products_price ON products(price);

-- Composite index for category + price filtering
CREATE INDEX IX_products_category_price ON products(category_id, price);

-- Index for slug lookups (SEO URLs)
CREATE UNIQUE INDEX IX_products_slug ON products(slug);

-- Index for created_at sorting (newest products)
CREATE INDEX IX_products_created_at ON products(created_at DESC);

-- ============================================
-- ORDER INDEXES
-- ============================================

-- Index on orders.user_id for user order history
CREATE INDEX IX_orders_user_id ON orders(user_id);

-- Index on orders.created_at for date-based queries
CREATE INDEX IX_orders_created_at ON orders(created_at DESC);

-- Composite index for user + status filtering
CREATE INDEX IX_orders_user_id_status ON orders(user_id, status);

-- Index for order status filtering
CREATE INDEX IX_orders_status ON orders(status);

-- ============================================
-- CATEGORY INDEXES
-- ============================================

-- Index for parent category lookups
CREATE INDEX IX_categories_parent_id ON categories(parent_id);

-- Index for slug lookups
CREATE UNIQUE INDEX IX_categories_slug ON categories(slug);

-- ============================================
-- CART INDEXES
-- ============================================

-- Index for user cart lookups
CREATE INDEX IX_cart_items_user_id ON cart_items(user_id);

-- Index for product in cart
CREATE INDEX IX_cart_items_product_id ON cart_items(product_id);

-- ============================================
-- REVIEW INDEXES
-- ============================================

-- Index for product reviews
CREATE INDEX IX_reviews_product_id ON reviews(product_id);

-- Index for user reviews
CREATE INDEX IX_reviews_user_id ON reviews(user_id);
```

### SQL Server Specific: Filtered Indexes

```sql
-- Filtered index for active products only (smaller, faster)
CREATE INDEX IX_products_active ON products(id, name, price, category_id)
WHERE status = 'active';

-- Filtered index for featured products
CREATE INDEX IX_products_featured_active ON products(id, name, price, image_url)
WHERE is_featured = 1 AND status = 'active';
```

### Analyze Index Usage

```sql
-- Check index usage statistics
SELECT 
    OBJECT_NAME(s.object_id) AS TableName,
    i.name AS IndexName,
    s.user_seeks,
    s.user_scans,
    s.user_lookups,
    s.user_updates
FROM sys.dm_db_index_usage_stats s
INNER JOIN sys.indexes i ON s.object_id = i.object_id AND s.index_id = i.index_id
WHERE OBJECT_NAME(s.object_id) IN ('products', 'orders', 'categories')
ORDER BY s.user_seeks DESC;
```

---

## PART 5 - GOLANG PAGINATION STRUCT

### PaginationQuery Struct

```go
package pkg/pagination

// PaginationQuery represents pagination parameters from API requests
type PaginationQuery struct {
    // Page is the page number (1-indexed, default: 1)
    Page int `form:"page" binding:"min=1"`
    
    // Limit is the number of items per page (default: 20, max: 100)
    Limit int `form:"limit" binding:"min=1,max=100"`
    
    // SortBy is the field to sort by
    SortBy string `form:"sort_by"`
    
    // SortOrder is the sort direction (asc, desc)
    SortOrder string `form:"sort_order"`
}

// Default returns a PaginationQuery with default values
func (p *PaginationQuery) Default() *PaginationQuery {
    if p.Page < 1 {
        p.Page = 1
    }
    if p.Limit < 1 || p.Limit > 100 {
        p.Limit = 20
    }
    if p.SortOrder == "" {
        p.SortOrder = "desc"
    }
    return p
}

// Offset calculates the offset based on page and limit
func (p *PaginationQuery) Offset() int {
    return (p.Page - 1) * p.Limit
}
```

### PaginationResult Struct

```go
package pkg/pagination

// PaginationResult represents pagination metadata in API responses
type PaginationResult struct {
    // Page is the current page number
    Page int `json:"page"`
    
    // Limit is the number of items per page
    Limit int `json:"limit"`
    
    // Total is the total number of items
    Total int64 `json:"total"`
    
    // TotalPages is the total number of pages
    TotalPages int `json:"total_pages"`
    
    // HasNext indicates if there is a next page
    HasNext bool `json:"has_next"`
    
    // HasPrev indicates if there is a previous page
    HasPrev bool `json:"has_prev"`
    
    // NextPage is the next page number (nil if none)
    NextPage *int `json:"next_page,omitempty"`
    
    // PrevPage is the previous page number (nil if none)
    PrevPage *int `json:"prev_page,omitempty"`
}

// NewPaginationResult creates a new PaginationResult
func NewPaginationResult(page, limit int, total int64) *PaginationResult {
    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }
    
    result := &PaginationResult{
        Page:       page,
        Limit:      limit,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }
    
    if result.HasNext {
        nextPage := page + 1
        result.NextPage = &nextPage
    }
    if result.HasPrev {
        prevPage := page - 1
        result.PrevPage = &prevPage
    }
    
    return result
}
```

### PaginatedResponse Generic Struct

```go
package pkg/response

import "ecommerce/pkg/pagination"

// PaginatedResponse represents a standard paginated API response
type PaginatedResponse[T any] struct {
    Success    bool                    `json:"success"`
    Message    string                  `json:"message,omitempty"`
    Data       []T                     `json:"data"`
    Pagination *pagination.PaginationResult `json:"pagination"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, pagination *pagination.PaginationResult) *PaginatedResponse[T] {
    return &PaginatedResponse[T]{
        Success:    true,
        Data:       data,
        Pagination: pagination,
    }
}
```

---

## PART 6 - PAGINATION HELPER

### Pagination Helper Functions

```go
package pkg/pagination

// Paginate calculates the offset for database queries
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

// CalculateTotalPages calculates total pages from total items and limit
func CalculateTotalPages(total int64, limit int) int {
    if limit <= 0 {
        limit = 20
    }
    pages := int(total) / limit
    if int(total)%limit > 0 {
        pages++
    }
    return pages
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

// ApplyPagination applies pagination to a GORM query
func ApplyPagination(query *gorm.DB, page, limit int) *gorm.DB {
    offset := Paginate(page, limit)
    return query.Offset(offset).Limit(limit)
}
```

### Usage Examples

```go
// Basic usage
offset := pagination.Paginate(3, 20)  // Returns 40

// In a repository method
func (r *productRepository) FindAll(limit, offset int) ([]Product, int64, error) {
    var products []Product
    var total int64
    
    // Count total
    r.db.Model(&Product{}).Where("status = ?", "active").Count(&total)
    
    // Apply pagination
    r.db.Offset(offset).Limit(limit).Find(&products)
    
    return products, total, nil
}

// Using GORM query builder
func GetProducts(db *gorm.DB, page, limit int) ([]Product, int64, error) {
    var products []Product
    var total int64
    
    query := db.Model(&Product{}).Where("status = ?", "active")
    query.Count(&total)
    
    pagination.ApplyPagination(query, page, limit).Find(&products)
    
    return products, total, nil
}
```

---

## PART 7 - PRODUCT LIST API

### API Endpoint

```http
GET /api/products?page=1&limit=20&sort_by=price&sort_order=asc
```

### Request Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number (1-indexed) |
| `limit` | int | 20 | Items per page (max: 100) |
| `sort_by` | string | created_at | Sort field (price, name, rating, sold_count) |
| `sort_order` | string | desc | Sort direction (asc, desc) |
| `category_id` | uint | - | Filter by category |
| `min_price` | float | - | Minimum price filter |
| `max_price` | float | - | Maximum price filter |
| `search` | string | - | Search keyword |

### Handler Implementation

```go
package handler

import (
    "ecommerce/pkg/pagination"
    "ecommerce/pkg/response"
    "github.com/gin-gonic/gin"
    "net/http"
)

// GetProducts handles GET /api/products
func (h *productHandler) GetProducts(c *gin.Context) {
    // Parse pagination parameters
    var query struct {
        Page       int     `form:"page" binding:"min=1"`
        Limit      int     `form:"limit" binding:"min=1,max=100"`
        SortBy     string  `form:"sort_by"`
        SortOrder  string  `form:"sort_order"`
        CategoryID *uint   `form:"category_id"`
        MinPrice   *float64 `form:"min_price"`
        MaxPrice   *float64 `form:"max_price"`
        Search     string  `form:"search"`
    }
    
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Apply defaults
    if query.Page < 1 {
        query.Page = 1
    }
    if query.Limit < 1 || query.Limit > 100 {
        query.Limit = 20
    }
    
    // Build filters
    filters := repository.ProductFilters{
        CategoryID: query.CategoryID,
        MinPrice:   query.MinPrice,
        MaxPrice:   query.MaxPrice,
        SortBy:     query.SortBy,
        SortOrder:  query.SortOrder,
    }
    
    // Get products from service
    products, total, err := h.productService.SearchProducts(query.Search, filters, query.Page, query.Limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }
    
    // Create pagination result
    paginationResult := pagination.NewPaginationResult(query.Page, query.Limit, total)
    
    // Return response
    c.JSON(http.StatusOK, response.NewPaginatedResponse(products, paginationResult))
}
```

### Example Response

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Wireless Bluetooth Headphones",
      "slug": "wireless-bluetooth-headphones",
      "price": 79.99,
      "original_price": 99.99,
      "discount_percent": 20,
      "rating_avg": 4.5,
      "rating_count": 128,
      "sold_count": 542,
      "stock": 150,
      "image_url": "/images/products/headphones.jpg",
      "category": {
        "id": 5,
        "name": "Electronics",
        "slug": "electronics"
      },
      "shop": {
        "id": 12,
        "name": "TechStore",
        "slug": "techstore"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 300,
    "total_pages": 15,
    "has_next": true,
    "has_prev": false,
    "next_page": 2,
    "prev_page": null
  }
}
```

---

## PART 8 - SEARCH OPTIMIZATION

### Optimizing Product Search

Efficient product search requires multiple optimization techniques:

### 1. Indexed Column Searches

```go
// Use indexed columns for filtering
func SearchProducts(db *gorm.DB, keyword string, categoryID *uint) ([]Product, int64, error) {
    query := db.Model(&Product{}).Where("status = ?", "active")
    
    // Use indexed column for category filter
    if categoryID != nil {
        query = query.Where("category_id = ?", *categoryID)  // Uses IX_products_category_id
    }
    
    // Full-text search on name (indexed)
    if keyword != "" {
        query = query.Where("name LIKE ?", "%"+keyword+"%")  // Uses IX_products_name
    }
    
    // ... execute query
}
```

### 2. Price Range Filtering

```go
// Efficient price range queries using index
func FilterByPrice(db *gorm.DB, minPrice, maxPrice float64) ([]Product, error) {
    var products []Product
    
    query := db.Model(&Product{}).Where("status = ?", "active")
    
    if minPrice > 0 {
        query = query.Where("price >= ?", minPrice)  // Uses IX_products_price
    }
    if maxPrice > 0 {
        query = query.Where("price <= ?", maxPrice)
    }
    
    query.Find(&products)
    return products, nil
}
```

### 3. Sorting Optimization

```go
// Use indexed columns for sorting
func GetSortedProducts(db *gorm.DB, sortBy, sortOrder string) ([]Product, error) {
    // Whitelist allowed sort columns
    allowedSorts := map[string]bool{
        "price": true, "created_at": true, 
        "sold_count": true, "rating_avg": true,
    }
    
    if !allowedSorts[sortBy] {
        sortBy = "created_at"
    }
    
    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }
    
    var products []Product
    db.Model(&Product{}).
        Where("status = ?", "active").
        Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).  // Uses appropriate index
        Find(&products)
    
    return products, nil
}
```

### 4. Combined Filter Query

```go
// Optimized combined filter with proper index usage
func SearchWithFilters(db *gorm.DB, filters ProductFilters) ([]Product, int64, error) {
    var products []Product
    var total int64
    
    // Start with most selective filter first (uses index)
    query := db.Model(&Product{}).Where("status = ?", "active")
    
    // Category filter (highly selective, uses index)
    if filters.CategoryID != nil {
        query = query.Where("category_id = ?", *filters.CategoryID)
    }
    
    // Price range (uses index)
    if filters.MinPrice != nil {
        query = query.Where("price >= ?", *filters.MinPrice)
    }
    if filters.MaxPrice != nil {
        query = query.Where("price <= ?", *filters.MaxPrice)
    }
    
    // Rating filter (uses index)
    if filters.MinRating != nil {
        query = query.Where("rating_avg >= ?", *filters.MinRating)
    }
    
    // Keyword search (least selective, apply last)
    if filters.Keyword != "" {
        query = query.Where("name LIKE ?", "%"+filters.Keyword+"%")
    }
    
    // Count total
    query.Count(&total)
    
    // Apply sorting (use indexed column)
    sortField := filters.SortBy
    if sortField == "" {
        sortField = "created_at"
    }
    query = query.Order(fmt.Sprintf("%s %s", sortField, filters.SortOrder))
    
    // Execute with pagination
    query.Offset(filters.Offset).Limit(filters.Limit).Find(&products)
    
    return products, total, nil
}
```

---

## PART 9 - DATABASE QUERY OPTIMIZATION

### Optimization Techniques

### 1. Select Specific Columns

```go
// ❌ Bad: Select all columns
var products []Product
db.Find(&products)

// ✅ Good: Select only needed columns
var products []ProductDTO
db.Model(&Product{}).
    Select("id, name, price, image_url, rating_avg, sold_count").
    Where("status = ?", "active").
    Find(&products)
```

### 2. Avoid N+1 Queries

```go
// ❌ Bad: N+1 query problem
var products []Product
db.Where("status = ?", "active").Find(&products)

for i := range products {
    // Separate query for each product's category
    db.First(&category, products[i].CategoryID)
    // Separate query for each product's images
    db.Where("product_id = ?", products[i].ID).Find(&images)
}

// ✅ Good: Eager loading with Preload
var products []Product
db.Preload("Category").
   Preload("Images", "is_primary = ?", true).
   Where("status = ?", "active").
   Find(&products)
```

### 3. Efficient JOINs

```go
// ✅ Use JOINs for filtering related data
var products []Product
db.Joins("JOIN categories ON products.category_id = categories.id").
   Where("categories.is_active = ? AND products.status = ?", true, "active").
   Preload("Category").
   Find(&products)
```

### 4. Use Subqueries for Aggregations

```go
// ✅ Efficient aggregation with subquery
var products []Product
db.Select("products.*, COUNT(orders.id) as order_count").
   Joins("LEFT JOIN order_items ON order_items.product_id = products.id").
   Joins("LEFT JOIN orders ON orders.id = order_items.order_id AND orders.status = ?", "completed").
   Group("products.id").
   Where("products.status = ?", "active").
   Find(&products)
```

### 5. Batch Operations

```go
// ❌ Bad: Individual inserts
for _, product := range products {
    db.Create(&product)
}

// ✅ Good: Batch insert
db.CreateInBatches(products, 100)

// ✅ Good: Bulk update
db.Model(&Product{}).
   Where("id IN ?", ids).
   Update("status", "active")
```

### Example Optimized Query

```go
func GetOptimizedProductList(db *gorm.DB, page, limit int, categoryID *uint) ([]Product, int64, error) {
    var products []Product
    var total int64
    
    // Build base query
    query := db.Model(&Product{}).Where("status = ?", "active")
    
    // Apply category filter (uses index)
    if categoryID != nil {
        query = query.Where("category_id = ?", *categoryID)
    }
    
    // Count total (separate query for accuracy)
    query.Count(&total)
    
    // Select specific columns + eager load
    err := query.
        Select("id, shop_id, category_id, name, slug, price, original_price, discount_percent, stock, sold_count, rating_avg, rating_count, status, is_featured").
        Preload("Category", "is_active = ?", true).
        Preload("Shop", "is_active = ?", true).
        Preload("Images", "is_primary = ?", true).
        Offset((page - 1) * limit).
        Limit(limit).
        Order("created_at DESC").
        Find(&products).Error
    
    return products, total, err
}
```

---

## PART 10 - SIMPLE CACHING

### Caching Strategies for E-Commerce

### 1. Cache Product List

```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type ProductCache struct {
    redis *redis.Client
    ttl   time.Duration
}

func NewProductCache(redis *redis.Client) *ProductCache {
    return &ProductCache{
        redis: redis,
        ttl:   10 * time.Minute,
    }
}

// CacheKey generates cache key for product list
func (c *ProductCache) CacheKey(page, limit int, categoryID *uint) string {
    if categoryID != nil {
        return fmt.Sprintf("products:cat%d:p%d:l%d", *categoryID, page, limit)
    }
    return fmt.Sprintf("products:all:p%d:l%d", page, limit)
}

// Get retrieves products from cache
func (c *ProductCache) Get(ctx context.Context, key string) ([]Product, int64, bool) {
    data, err := c.redis.Get(ctx, key).Bytes()
    if err != nil {
        return nil, 0, false
    }
    
    var cached struct {
        Products []Product `json:"products"`
        Total    int64     `json:"total"`
    }
    
    if err := json.Unmarshal(data, &cached); err != nil {
        return nil, 0, false
    }
    
    return cached.Products, cached.Total, true
}

// Set stores products in cache
func (c *ProductCache) Set(ctx context.Context, key string, products []Product, total int64) error {
    data, _ := json.Marshal(map[string]interface{}{
        "products": products,
        "total":    total,
    })
    
    return c.redis.Set(ctx, key, data, c.ttl).Err()
}

// Invalidate clears cache for product list
func (c *ProductCache) Invalidate(ctx context.Context, categoryID *uint) error {
    // Invalidate all product list caches when product is updated
    pattern := "products:*"
    iter := c.redis.Scan(ctx, 0, pattern, 100).Iterator()
    
    for iter.Next(ctx) {
        c.redis.Del(ctx, iter.Val())
    }
    
    return iter.Err()
}
```

### 2. Cache Product Details

```go
// Cache single product details
func (c *ProductCache) GetProduct(ctx context.Context, id uint) (*Product, bool) {
    key := fmt.Sprintf("product:detail:%d", id)
    
    data, err := c.redis.Get(ctx, key).Bytes()
    if err != nil {
        return nil, false
    }
    
    var product Product
    if err := json.Unmarshal(data, &product); err != nil {
        return nil, false
    }
    
    return &product, true
}

func (c *ProductCache) SetProduct(ctx context.Context, product *Product) error {
    key := fmt.Sprintf("product:detail:%d", product.ID)
    data, _ := json.Marshal(product)
    
    // Longer TTL for product details
    return c.redis.Set(ctx, key, data, 30*time.Minute).Err()
}

func (c *ProductCache) DeleteProduct(ctx context.Context, id uint) error {
    key := fmt.Sprintf("product:detail:%d", id)
    return c.redis.Del(ctx, key).Err()
}
```

### 3. When to Use Caching

| Scenario | Cache Strategy | TTL |
|----------|---------------|-----|
| Product list (homepage) | Cache with short TTL | 5-10 min |
| Product details | Cache with medium TTL | 30 min |
| Category tree | Cache with long TTL | 1 hour |
| User-specific data | No cache or very short TTL | 1-5 min |
| Flash sale products | Cache with very short TTL | 1-2 min |
| Search results | Cache with medium TTL | 15 min |

### 4. Cache Invalidation

```go
// Invalidate cache when product is updated
func (s *productService) UpdateProduct(id uint, product *Product) error {
    // Update database
    err := s.repo.Update(product)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    ctx := context.Background()
    s.productCache.DeleteProduct(ctx, id)
    s.productCache.Invalidate(ctx, &product.CategoryID)
    
    return nil
}
```

---

## PART 11 - LOAD HANDLING

### Handling 50+ Concurrent Users

### 1. Connection Pooling

```go
// Configure database connection pool
func InitDatabase(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Get underlying SQL DB
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    sqlDB.SetMaxIdleConns(25)        // Max idle connections
    sqlDB.SetMaxOpenConns(100)       // Max open connections
    sqlDB.SetConnMaxLifetime(5 * time.Minute)  // Connection lifetime
    
    return db, nil
}
```

### 2. Efficient Queries

```go
// With 50 concurrent users, each requesting 20 products:
// Without optimization: 50 × 500,000 rows = 25 million rows
// With pagination: 50 × 20 rows = 1,000 rows

// Optimized query for concurrent load
func GetProductsUnderLoad(db *gorm.DB, page, limit int) ([]Product, int64, error) {
    // Use read replica if available (reduces primary DB load)
    // Use covering index (select only indexed columns)
    
    var products []Product
    var total int64
    
    // Use context with timeout to prevent long-running queries
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    query := db.WithContext(ctx).Model(&Product{}).Where("status = ?", "active")
    query.Count(&total)
    
    // Select only necessary columns (covering index)
    query.Select("id, name, price, image_url, rating_avg, sold_count, category_id").
        Limit(limit).
        Offset((page - 1) * limit).
        Order("created_at DESC").
        Find(&products)
    
    return products, total, nil
}
```

### 3. Rate Limiting

```go
// Implement rate limiting middleware
func RateLimiter() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(time.Second), 100) // 100 req/sec
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Too many requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 4. Load Distribution

```
┌─────────────────────────────────────────────────────────┐
│                    Load Balancer                         │
│                    (Nginx / HAProxy)                     │
└─────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│   App Server  │    │   App Server  │    │   App Server  │
│   Instance 1  │    │   Instance 2  │    │   Instance 3  │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  Primary DB   │    │  Read Replica │    │     Redis     │
│  (SQL Server) │    │  (SQL Server) │    │    Cache      │
└───────────────┘    └───────────────┘    └───────────────┘
```

---

## PART 12 - EXAMPLE IMPLEMENTATION

### Complete Implementation Example

### 1. Pagination Helper (pkg/pagination/pagination.go)

```go
package pagination

import (
    "gorm.io/gorm"
    "math"
)

// Query represents pagination query parameters
type Query struct {
    Page      int    `form:"page" binding:"min=1"`
    Limit     int    `form:"limit" binding:"min=1,max=100"`
    SortBy    string `form:"sort_by"`
    SortOrder string `form:"sort_order"`
}

// Result represents pagination metadata
type Result struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
    NextPage   *int  `json:"next_page,omitempty"`
    PrevPage   *int  `json:"prev_page,omitempty"`
}

// NewQuery creates a new Query with defaults
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

// Offset calculates the database offset
func (q *Query) Offset() int {
    return (q.Page - 1) * q.Limit
}

// Apply applies pagination to a GORM query
func (q *Query) Apply(query *gorm.DB) *gorm.DB {
    return query.Offset(q.Offset()).Limit(q.Limit)
}

// NewResult creates a new pagination Result
func NewResult(page, limit int, total int64) *Result {
    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    
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
```

### 2. Product List API (api/products/products.go)

```go
package products

import (
    "ecommerce/internal/handler"
    "ecommerce/pkg/pagination"
    "ecommerce/pkg/response"
    "github.com/gin-gonic/gin"
    "net/http"
)

// SetupRoutes configures product routes
func SetupRoutes(rg *gin.RouterGroup, h *handler.ProductHandler) {
    products := rg.Group("/products")
    {
        products.GET("", h.GetProducts)
        products.GET("/:id", h.GetProductByID)
        products.POST("", h.CreateProduct)
        products.PUT("/:id", h.UpdateProduct)
        products.DELETE("/:id", h.DeleteProduct)
    }
}
```

### 3. Handler Implementation (internal/handler/product_handler.go)

```go
package handler

import (
    "ecommerce/internal/domain/model"
    "ecommerce/internal/service"
    "ecommerce/pkg/pagination"
    "ecommerce/pkg/response"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

type ProductHandler struct {
    productService service.ProductServiceEnhanced
}

func NewProductHandler(ps service.ProductServiceEnhanced) *ProductHandler {
    return &ProductHandler{productService: ps}
}

// GetProducts handles GET /api/products
func (h *ProductHandler) GetProducts(c *gin.Context) {
    // Parse query parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    
    // Validate
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    // Build filters
    filters := service.ProductFilters{
        SortBy:    c.DefaultQuery("sort_by", "created_at"),
        SortOrder: c.DefaultQuery("sort_order", "desc"),
    }
    
    // Parse optional filters
    if categoryID := c.Query("category_id"); categoryID != "" {
        if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
            catID := uint(id)
            filters.CategoryID = &catID
        }
    }
    
    if minPrice := c.Query("min_price"); minPrice != "" {
        if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
            filters.MinPrice = &price
        }
    }
    
    if maxPrice := c.Query("max_price"); maxPrice != "" {
        if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
            filters.MaxPrice = &price
        }
    }
    
    // Get products
    products, total, err := h.productService.SearchProducts(c.Query("search"), filters, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }
    
    // Create pagination result
    paginationResult := pagination.NewResult(page, limit, total)
    
    // Return response
    c.JSON(http.StatusOK, response.PaginatedResponse[model.Product]{
        Success:    true,
        Message:    "Products retrieved successfully",
        Data:       products,
        Pagination: paginationResult,
    })
}
```

### 4. Optimized Repository Query (internal/repository/product_repository.go)

```go
package repository

import (
    "ecommerce/internal/domain/model"
    "fmt"
    "strings"
    "gorm.io/gorm"
)

// FindAllOptimized returns paginated products with optimized query
func (r *productRepositoryEnhanced) FindAllOptimized(page, limit int) ([]model.Product, int64, error) {
    var products []model.Product
    var total int64
    
    // Use optimized query with specific columns
    query := r.db.Model(&model.Product{}).
        Where("status = ?", model.ProductStatusActive)
    
    // Count total
    query.Count(&total)
    
    // Select specific columns + eager load
    err := query.
        Select("id, shop_id, category_id, name, slug, description, price, original_price, discount_percent, stock, sold_count, rating_avg, rating_count, status, is_featured, created_at").
        Preload("Category", "is_active = ?", true).
        Preload("Shop", "is_active = ?", true).
        Preload("Images", "is_primary = ?", true).
        Offset((page - 1) * limit).
        Limit(limit).
        Order("created_at DESC").
        Find(&products).Error
    
    return products, total, err
}

// SearchOptimized performs optimized search with filters
func (r *productRepositoryEnhanced) SearchOptimized(keyword string, filters ProductFilters, page, limit int) ([]model.Product, int64, error) {
    var products []model.Product
    var total int64
    
    query := r.db.Model(&model.Product{}).Where("status = ?", model.ProductStatusActive)
    
    // Apply most selective filters first
    if filters.CategoryID != nil && *filters.CategoryID > 0 {
        query = query.Where("category_id = ?", *filters.CategoryID)
    }
    
    if filters.MinPrice != nil && *filters.MinPrice > 0 {
        query = query.Where("price >= ?", *filters.MinPrice)
    }
    
    if filters.MaxPrice != nil && *filters.MaxPrice > 0 {
        query = query.Where("price <= ?", *filters.MaxPrice)
    }
    
    if filters.MinRating != nil && *filters.MinRating > 0 {
        query = query.Where("rating_avg >= ?", *filters.MinRating)
    }
    
    // Keyword search last (least selective)
    if keyword != "" {
        searchPattern := "%" + strings.TrimSpace(keyword) + "%"
        query = query.Where("name LIKE ? OR short_description LIKE ? OR brand LIKE ?",
            searchPattern, searchPattern, searchPattern)
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Validate sort field
    validSortFields := map[string]bool{
        "price": true, "rating_avg": true, "sold_count": true,
        "created_at": true, "view_count": true, "name": true,
    }
    sortBy := filters.SortBy
    if !validSortFields[sortBy] {
        sortBy = "created_at"
    }
    sortOrder := strings.ToUpper(filters.SortOrder)
    if sortOrder != "ASC" && sortOrder != "DESC" {
        sortOrder = "DESC"
    }
    
    // Execute query
    err := query.
        Select("id, shop_id, category_id, name, slug, price, original_price, discount_percent, stock, sold_count, rating_avg, rating_count, status, is_featured, created_at").
        Preload("Category", "is_active = ?", true).
        Preload("Shop", "is_active = ?", true).
        Preload("Images", "is_primary = ?", true).
        Offset((page - 1) * limit).
        Limit(limit).
        Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
        Find(&products).Error
    
    return products, total, err
}
```

---

## Summary

This implementation provides:

1. **Pagination**: Efficient data retrieval with configurable page size
2. **Database Indexing**: Optimized indexes for common queries
3. **Query Optimization**: Selective columns, eager loading, efficient JOINs
4. **Caching**: Redis-based caching for frequently accessed data
5. **Load Handling**: Connection pooling, rate limiting, efficient queries

The system can now handle 50+ concurrent users efficiently while maintaining fast response times and preventing database overload.
