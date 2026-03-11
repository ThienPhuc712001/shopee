# Product Search and Filter System - Documentation

## Overview

The Product Search and Filter System provides powerful search capabilities with multiple filters, sorting options, and pagination for the e-commerce platform.

---

## Table of Contents

1. [Business Flow](#business-flow)
2. [Search Parameters](#search-parameters)
3. [Database Optimization](#database-optimization)
4. [API Endpoints](#api-endpoints)
5. [Frontend Integration](#frontend-integration)
6. [Examples](#examples)

---

## Business Flow

### How Product Search Works

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  1. User    │────▶│  2. Frontend │────▶│  3. Backend │
│  Enters     │     │  Sends       │     │  Queries    │
│  Keyword    │     │  Request     │     │  Database   │
└─────────────┘     └──────────────┘     └─────────────┘
                                              │
                                              ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  5. Backend │◀────│  4. Results  │◀────│  Filters &  │
│  Returns    │     │  Paginated   │     │  Sorts      │
│  Results    │     │              │     │             │
└─────────────┘     └──────────────┘     └─────────────┘
```

### Step-by-Step Flow

| Step | Actor | Action | Result |
|------|-------|--------|--------|
| 1 | User | Enters keyword in search bar | Search query formed |
| 2 | Frontend | Sends GET request with filters | Query parameters sent |
| 3 | Backend | Queries database with filters | SQL query executed |
| 4 | Backend | Filters and sorts results | Results ordered and paginated |
| 5 | Backend | Returns paginated response | JSON response with products |

---

## Search Parameters

### Supported Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `keyword` | string | Search in name, description, brand, SKU | `iphone` |
| `category_id` | uint | Filter by category | `5` |
| `shop_id` | uint | Filter by seller shop | `10` |
| `min_price` | float64 | Minimum price filter | `500` |
| `max_price` | float64 | Maximum price filter | `2000` |
| `min_rating` | float64 | Minimum rating (0-5) | `4.5` |
| `brands` | string | Comma-separated brands | `Apple,Samsung` |
| `sort` | string | Sorting option | `price_asc` |
| `page` | int | Page number | `1` |
| `limit` | int | Items per page (max 100) | `20` |

### Sort Options

| Value | Description |
|-------|-------------|
| `newest` | Newest products first |
| `price_asc` | Price: Low to High |
| `price_desc` | Price: High to Low |
| `best_selling` | Best selling products |
| `top_rated` | Highest rated products |

---

## Database Optimization

### Indexes on Products Table

```sql
-- Name search index
CREATE INDEX IX_products_name ON products(name);

-- Category filter index
CREATE INDEX IX_products_category_id ON products(category_id);

-- Price filter index
CREATE INDEX IX_products_price ON products(price);

-- Rating filter index
CREATE INDEX IX_products_rating_avg ON products(rating_avg);

-- Status filter index
CREATE INDEX IX_products_status ON products(status);

-- Sold count for best sellers
CREATE INDEX IX_products_sold_count ON products(sold_count);

-- Composite index for common queries
CREATE INDEX IX_products_category_status_price 
ON products(category_id, status, price);
```

### Why Indexing Improves Performance

| Index | Benefit |
|-------|---------|
| `name` | Faster LIKE searches |
| `category_id` | Quick category filtering |
| `price` | Efficient price range queries |
| `rating_avg` | Fast rating filtering |
| `status` | Quick active product filtering |
| `sold_count` | Fast best-seller queries |

---

## API Endpoints

### 1. Search Products

```http
GET /api/products/search
```

**Request Parameters:**
```
?keyword=iphone&category_id=5&min_price=500&max_price=2000&sort=price_asc&page=1&limit=20
```

**Response:**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": 1,
        "name": "iPhone 15 Pro",
        "slug": "iphone-15-pro",
        "price": 999,
        "original_price": 1099,
        "discount_percent": 9,
        "rating_avg": 4.8,
        "rating_count": 156,
        "sold_count": 523,
        "stock": 50,
        "available_stock": 45,
        "image_url": "/uploads/products/iphone15.jpg",
        "category": {
          "id": 5,
          "name": "Phones"
        },
        "shop": {
          "id": 10,
          "name": "Apple Store"
        }
      }
    ],
    "total": 120,
    "page": 1,
    "limit": 20,
    "total_pages": 6
  }
}
```

### 2. Get Product by ID

```http
GET /api/products/:id
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "description": "...",
    "price": 999,
    "images": [...],
    "category": {...},
    "shop": {...}
  }
}
```

### 3. Get Featured Products

```http
GET /api/products/featured?limit=8
```

### 4. Get Flash Sale Products

```http
GET /api/products/flash-sale
```

### 5. Get Best Sellers

```http
GET /api/products/best-sellers?limit=8
```

### 6. Get Related Products

```http
GET /api/products/:id/related?limit=8
```

### 7. Get Search Suggestions

```http
GET /api/products/suggest?q=iph&limit=5
```

**Response:**
```json
{
  "success": true,
  "data": {
    "suggestions": [
      {"id": 1, "name": "iPhone 15 Pro", "slug": "iphone-15-pro", "price": 999},
      {"id": 2, "name": "iPhone 14", "slug": "iphone-14", "price": 799}
    ],
    "count": 2
  }
}
```

### 8. Get Search Filters

```http
GET /api/products/search/filters?category_id=5
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sort_options": [
      {"value": "newest", "label": "Newest"},
      {"value": "price_asc", "label": "Price: Low to High"},
      {"value": "price_desc", "label": "Price: High to Low"},
      {"value": "best_selling", "label": "Best Selling"},
      {"value": "top_rated", "label": "Top Rated"}
    ],
    "price_range": {"min": 0, "max": 10000},
    "rating_range": {"min": 0, "max": 5}
  }
}
```

---

## Frontend Integration

### React Search Component

```jsx
function ProductSearch() {
  const [products, setProducts] = useState([]);
  const [filters, setFilters] = useState({
    keyword: '',
    category_id: null,
    min_price: null,
    max_price: null,
    min_rating: null,
    sort: 'newest',
    page: 1,
    limit: 20
  });
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({});

  const searchProducts = async () => {
    setLoading(true);
    
    // Build query string
    const params = new URLSearchParams();
    if (filters.keyword) params.append('keyword', filters.keyword);
    if (filters.category_id) params.append('category_id', filters.category_id);
    if (filters.min_price) params.append('min_price', filters.min_price);
    if (filters.max_price) params.append('max_price', filters.max_price);
    if (filters.min_rating) params.append('min_rating', filters.min_rating);
    if (filters.sort) params.append('sort', filters.sort);
    params.append('page', filters.page);
    params.append('limit', filters.limit);

    const response = await fetch(`/api/products/search?${params}`);
    const result = await response.json();

    if (result.success) {
      setProducts(result.data.products);
      setPagination({
        total: result.data.total,
        page: result.data.page,
        limit: result.data.limit,
        total_pages: result.data.total_pages
      });
    }
    setLoading(false);
  };

  useEffect(() => {
    searchProducts();
  }, [filters]);

  return (
    <div>
      {/* Search Bar */}
      <input
        type="text"
        placeholder="Search products..."
        value={filters.keyword}
        onChange={(e) => setFilters({...filters, keyword: e.target.value})}
      />

      {/* Category Filter */}
      <select
        value={filters.category_id || ''}
        onChange={(e) => setFilters({...filters, category_id: e.target.value ? Number(e.target.value) : null})}
      >
        <option value="">All Categories</option>
        {/* Category options */}
      </select>

      {/* Price Range */}
      <input
        type="number"
        placeholder="Min Price"
        value={filters.min_price || ''}
        onChange={(e) => setFilters({...filters, min_price: e.target.value ? Number(e.target.value) : null})}
      />
      <input
        type="number"
        placeholder="Max Price"
        value={filters.max_price || ''}
        onChange={(e) => setFilters({...filters, max_price: e.target.value ? Number(e.target.value) : null})}
      />

      {/* Sort Dropdown */}
      <select
        value={filters.sort}
        onChange={(e) => setFilters({...filters, sort: e.target.value})}
      >
        <option value="newest">Newest</option>
        <option value="price_asc">Price: Low to High</option>
        <option value="price_desc">Price: High to Low</option>
        <option value="best_selling">Best Selling</option>
        <option value="top_rated">Top Rated</option>
      </select>

      {/* Products Grid */}
      {loading ? (
        <LoadingSpinner />
      ) : (
        <div className="products-grid">
          {products.map(product => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}

      {/* Pagination */}
      <Pagination
        currentPage={pagination.page}
        totalPages={pagination.total_pages}
        onPageChange={(page) => setFilters({...filters, page})}
      />
    </div>
  );
}
```

### Search Bar with Debounce

```jsx
function SearchBar({ onSearch }) {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      if (query.length >= 2) {
        fetch(`/api/products/suggest?q=${query}&limit=5`)
          .then(res => res.json())
          .then(data => {
            if (data.success) {
              setSuggestions(data.data.suggestions);
            }
          });
      } else {
        setSuggestions([]);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [query]);

  return (
    <div className="search-bar">
      <input
        type="text"
        placeholder="Search products..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onSubmit={() => onSearch(query)}
      />
      
      {suggestions.length > 0 && (
        <div className="suggestions">
          {suggestions.map(item => (
            <Link key={item.id} to={`/products/${item.slug}`}>
              {item.name} - ${item.price}
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
```

### Price Range Slider

```jsx
function PriceRangeFilter({ onFilterChange }) {
  const [range, setRange] = useState([0, 10000]);

  return (
    <div className="price-filter">
      <h3>Price Range</h3>
      <Slider
        range
        min={0}
        max={10000}
        value={range}
        onChange={(value) => {
          setRange(value);
          onFilterChange({
            min_price: value[0],
            max_price: value[1]
          });
        }}
      />
      <div className="price-values">
        <span>${range[0]}</span>
        <span>${range[1]}</span>
      </div>
    </div>
  );
}
```

---

## Examples

### cURL Examples

**Search with filters:**
```bash
curl -X GET "http://localhost:8080/api/products/search?keyword=iphone&min_price=500&max_price=2000&sort=price_asc&page=1&limit=20"
```

**Get featured products:**
```bash
curl -X GET "http://localhost:8080/api/products/featured?limit=8"
```

**Get best sellers:**
```bash
curl -X GET "http://localhost:8080/api/products/best-sellers?limit=8"
```

**Get search suggestions:**
```bash
curl -X GET "http://localhost:8080/api/products/suggest?q=iph&limit=5"
```

**Get related products:**
```bash
curl -X GET "http://localhost:8080/api/products/1/related?limit=8"
```

---

## Pagination

### How OFFSET and LIMIT Work

```sql
-- Page 1, Limit 20
SELECT * FROM products
ORDER BY created_at DESC
OFFSET 0 ROWS
FETCH NEXT 20 ROWS ONLY;

-- Page 2, Limit 20
SELECT * FROM products
ORDER BY created_at DESC
OFFSET 20 ROWS
FETCH NEXT 20 ROWS ONLY;

-- Page 3, Limit 20
SELECT * FROM products
ORDER BY created_at DESC
OFFSET 40 ROWS
FETCH NEXT 20 ROWS ONLY;
```

### Formula

```
OFFSET = (page - 1) * limit
```

### Response Structure

```json
{
  "success": true,
  "data": {
    "products": [...],
    "total": 120,       // Total products matching filters
    "page": 1,          // Current page
    "limit": 20,        // Items per page
    "total_pages": 6    // Total pages (ceil(total / limit))
  }
}
```

---

## Search Query Logic

### SQL Query Conditions

```sql
SELECT * FROM products
WHERE deleted_at IS NULL
  AND status = 'active'
  AND (
    LOWER(name) LIKE '%keyword%'
    OR LOWER(description) LIKE '%keyword%'
    OR LOWER(brand) LIKE '%keyword%'
    OR LOWER(sku) LIKE '%keyword%'
  )
  AND (category_id = ? OR ? IS NULL)
  AND (price >= ? OR ? IS NULL)
  AND (price <= ? OR ? IS NULL)
  AND (rating_avg >= ? OR ? IS NULL)
ORDER BY created_at DESC
OFFSET 0 ROWS
FETCH NEXT 20 ROWS ONLY;
```

### Filter Combination

All filters are combined with **AND** logic:
- Keyword search (name, description, brand, SKU)
- **AND** Category filter
- **AND** Price range
- **AND** Rating filter
- **AND** Status filter (default: active)

---

## Implementation Summary

### Files Used

| File | Purpose |
|------|---------|
| `internal/repository/product_repository_enhanced.go` | Database queries |
| `internal/service/product_service_enhanced.go` | Business logic |
| `internal/handler/product_handler_enhanced.go` | HTTP handlers |
| `api/routes_enhanced.go` | Route configuration |

### Key Functions

**Repository:**
- `Search(keyword, filters, limit, offset)` - Search with filters
- `FindFeatured(limit)` - Get featured products
- `FindBestSellers(limit)` - Get best sellers

**Service:**
- `SearchProducts(keyword, filters, page, limit)` - Search products
- `GetFeaturedProducts(limit)` - Get featured products
- `GetBestSellers(limit)` - Get best sellers

**Handler:**
- `SearchProducts(c)` - HTTP search endpoint
- `GetFeaturedProducts(c)` - HTTP featured endpoint
- `GetBestSellers(c)` - HTTP best sellers endpoint

---

## Best Practices

1. **Debounce search input** - Wait 300ms before sending request
2. **Cache category list** - Categories don't change often
3. **Use pagination** - Never load all products at once
4. **Show loading states** - Indicate when searching
5. **Preserve filters** - Keep filters when navigating pages
6. **Show result count** - Display total products found
7. **Support URL sharing** - Encode filters in URL params

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Slow search | Add indexes on name, category_id, price |
| No results | Check if products are active and not deleted |
| Pagination broken | Verify OFFSET calculation |
| Sort not working | Check sort parameter values |

---

## Performance Tips

1. **Use database indexes** - Critical for large datasets
2. **Limit result count** - Max 100 items per page
3. **Cache popular searches** - Redis for frequent queries
4. **Preload associations** - Category, Shop, Images
5. **Use read replicas** - For high-traffic sites

---

✅ **System Ready**

The Product Search and Filter System is fully implemented with:
- Multi-field keyword search
- Price range filtering
- Rating filtering
- Category filtering
- Multiple sort options
- Efficient pagination
- Search suggestions
- Related products
