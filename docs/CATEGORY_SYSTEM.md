# Product Categories System - Documentation

## Overview

The Categories System provides hierarchical product categorization for the e-commerce platform, allowing admins to organize products and users to browse by category.

---

## Table of Contents

1. [Business Flow](#business-flow)
2. [Category Structure](#category-structure)
3. [Database Schema](#database-schema)
4. [API Endpoints](#api-endpoints)
5. [Frontend Usage](#frontend-usage)
6. [Examples](#examples)

---

## Business Flow

### How Categories Work

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  1. Admin   │────▶│  2. Products │────▶│  3. Users   │
│  Creates    │     │  Assigned to │     │  Browse     │
│  Categories │     │  Categories  │     │  by Category│
└─────────────┘     └──────────────┘     └─────────────┘
```

### Step-by-Step Flow

| Step | Actor | Action | Result |
|------|-------|--------|--------|
| 1 | Admin | Creates categories | Category tree structure |
| 2 | Admin/Seller | Assigns products to categories | Products organized |
| 3 | User | Browses categories | Sees category navigation |
| 4 | User | Clicks category | Views products in category |

---

## Category Structure

### Hierarchical Categories

Supports unlimited nesting levels:

```
Electronics (Level 0)
├── Phones (Level 1)
│   ├── Smartphones (Level 2)
│   └── Feature Phones (Level 2)
├── Laptops (Level 1)
│   ├── Gaming (Level 2)
│   └── Business (Level 2)
└── Tablets (Level 1)

Fashion (Level 0)
├── Men (Level 1)
│   ├── Clothing (Level 2)
│   └── Shoes (Level 2)
├── Women (Level 1)
│   ├── Clothing (Level 2)
│   └── Shoes (Level 2)
└── Kids (Level 1)
```

### Category Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | uint | Unique identifier |
| `parent_id` | *uint | Parent category (NULL for root) |
| `name` | string | Category name |
| `slug` | string | URL-friendly identifier |
| `description` | string | Category description |
| `icon_url` | string | Icon image URL |
| `image_url` | string | Banner image URL |
| `level` | int | Depth in hierarchy (0 = root) |
| `sort_order` | int | Display order |
| `is_active` | bool | Visibility status |
| `attributes` | string | JSON custom attributes |

---

## Database Schema

### Categories Table

```sql
CREATE TABLE categories (
    id INT IDENTITY(1,1) PRIMARY KEY,
    parent_id INT,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    image_url VARCHAR(500),
    level INT DEFAULT 0,
    sort_order INT DEFAULT 0,
    is_active BIT DEFAULT 1,
    attributes TEXT,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    deleted_at DATETIME,
    
    CONSTRAINT FK_categories_parent 
        FOREIGN KEY (parent_id) REFERENCES categories(id),
    
    INDEX IX_categories_parent_id (parent_id),
    INDEX IX_categories_slug (slug),
    INDEX IX_categories_level (level),
    INDEX IX_categories_sort_order (sort_order)
);
```

### Products Relationship

```sql
-- Products table has category_id foreign key
ALTER TABLE products ADD
    category_id INT NOT NULL,
    CONSTRAINT FK_products_category 
        FOREIGN KEY (category_id) REFERENCES categories(id);

CREATE INDEX IX_products_category_id ON products(category_id);
```

---

## API Endpoints

### Public Endpoints

#### 1. Get All Categories

```http
GET /api/categories
```

**Response:**
```json
{
  "success": true,
  "data": {
    "categories": [
      {
        "id": 1,
        "parent_id": null,
        "name": "Electronics",
        "slug": "electronics",
        "description": "Electronic devices",
        "icon_url": "/icons/electronics.png",
        "image_url": "/images/electronics.jpg",
        "level": 0,
        "sort_order": 1,
        "is_active": true,
        "created_at": "2026-03-10T10:00:00Z",
        "updated_at": "2026-03-10T10:00:00Z"
      }
    ],
    "count": 5
  }
}
```

#### 2. Get Category Tree

```http
GET /api/categories/tree
```

**Response:**
```json
{
  "success": true,
  "data": {
    "tree": [
      {
        "id": 1,
        "name": "Electronics",
        "slug": "electronics",
        "level": 0,
        "product_count": 150,
        "children": [
          {
            "id": 2,
            "name": "Phones",
            "slug": "phones",
            "level": 1,
            "product_count": 50,
            "children": []
          }
        ]
      }
    ]
  }
}
```

#### 3. Get Products by Category

```http
GET /api/categories/:id/products?page=1&limit=20
```

**Response:**
```json
{
  "success": true,
  "data": {
    "products": [...],
    "pagination": {
      "total": 150,
      "page": 1,
      "limit": 20,
      "total_pages": 8
    }
  }
}
```

#### 4. Get Category Breadcrumb

```http
GET /api/categories/:id/breadcrumb
```

**Response:**
```json
{
  "success": true,
  "data": {
    "breadcrumb": [
      {"id": 1, "name": "Electronics", "slug": "electronics"},
      {"id": 2, "name": "Phones", "slug": "phones"},
      {"id": 5, "name": "Smartphones", "slug": "smartphones"}
    ]
  }
}
```

#### 5. Get Featured Categories

```http
GET /api/categories/featured?limit=8
```

#### 6. Search Categories

```http
GET /api/categories/search?q=phone
```

### Admin Endpoints (Require Authentication)

#### 1. Create Category

```http
POST /api/categories
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "name": "Smartphones",
  "parent_id": 2,
  "description": "Latest smartphones",
  "icon_url": "/icons/smartphones.png",
  "image_url": "/images/smartphones.jpg",
  "sort_order": 1,
  "is_active": true
}
```

**Response:**
```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": 5,
    "name": "Smartphones",
    "slug": "smartphones-abc123",
    "parent_id": 2,
    "level": 2,
    ...
  }
}
```

#### 2. Update Category

```http
PUT /api/categories/:id
Authorization: Bearer {token}
```

#### 3. Delete Category

```http
DELETE /api/categories/:id
Authorization: Bearer {token}
```

**Note:** Cannot delete if category has subcategories or products.

---

## Frontend Usage

### Navigation Menu

```jsx
function CategoryNavigation() {
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    fetch('/api/categories/tree')
      .then(res => res.json())
      .then(data => setCategories(data.data.tree));
  }, []);

  return (
    <nav>
      {categories.map(cat => (
        <div key={cat.id}>
          <Link to={`/categories/${cat.slug}`}>{cat.name}</Link>
          {cat.children.length > 0 && (
            <ul>
              {cat.children.map(child => (
                <li key={child.id}>
                  <Link to={`/categories/${child.slug}`}>{child.name}</Link>
                </li>
              ))}
            </ul>
          )}
        </div>
      ))}
    </nav>
  );
}
```

### Product Filter by Category

```jsx
function ProductFilter({ onSelectCategory }) {
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    fetch('/api/categories')
      .then(res => res.json())
      .then(data => setCategories(data.data.categories));
  }, []);

  return (
    <div>
      <h3>Categories</h3>
      {categories.map(cat => (
        <button 
          key={cat.id}
          onClick={() => onSelectCategory(cat.id)}
        >
          {cat.name} ({cat.product_count})
        </button>
      ))}
    </div>
  );
}
```

### Homepage Category Section

```jsx
function FeaturedCategories() {
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    fetch('/api/categories/featured?limit=8')
      .then(res => res.json())
      .then(data => setCategories(data.data.categories));
  }, []);

  return (
    <section>
      <h2>Shop by Category</h2>
      <div className="grid">
        {categories.map(cat => (
          <Link 
            key={cat.id} 
            to={`/categories/${cat.slug}`}
            className="card"
          >
            <img src={cat.image_url} alt={cat.name} />
            <h3>{cat.name}</h3>
            <p>{cat.product_count} products</p>
          </Link>
        ))}
      </div>
    </section>
  );
}
```

### Breadcrumb Component

```jsx
function Breadcrumb({ categoryId }) {
  const [breadcrumb, setBreadcrumb] = useState([]);

  useEffect(() => {
    fetch(`/api/categories/${categoryId}/breadcrumb`)
      .then(res => res.json())
      .then(data => setBreadcrumb(data.data.breadcrumb));
  }, [categoryId]);

  return (
    <nav aria-label="breadcrumb">
      {breadcrumb.map((item, index) => (
        <span key={item.id}>
          <Link to={`/categories/${item.slug}`}>{item.name}</Link>
          {index < breadcrumb.length - 1 && ' > '}
        </span>
      ))}
    </nav>
  );
}
```

---

## Examples

### cURL Examples

**Create Category:**
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptops",
    "parent_id": 3,
    "description": "High-performance gaming laptops",
    "is_active": true
  }'
```

**Get Category Tree:**
```bash
curl -X GET http://localhost:8080/api/categories/tree
```

**Get Products by Category:**
```bash
curl -X GET "http://localhost:8080/api/categories/5/products?page=1&limit=20"
```

**Delete Category:**
```bash
curl -X DELETE http://localhost:8080/api/categories/5 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## Validation Rules

| Rule | Description | Error |
|------|-------------|-------|
| Name required | Category must have a name | "Category name is required" |
| Name length | 2-255 characters | "Category name must be between 2 and 255 characters" |
| Parent exists | Parent category must exist | "Parent category not found" |
| No circular refs | Cannot be child of own descendant | "Circular category hierarchy detected" |
| No delete with children | Must delete subcategories first | "Cannot delete category with subcategories" |
| No delete with products | Must move products first | "Cannot delete category with products" |

---

## Best Practices

1. **Keep hierarchy shallow** - Max 3-4 levels for better UX
2. **Use descriptive names** - Clear category names help navigation
3. **Add product counts** - Show users how many products in each category
4. **Optimize images** - Compress category images for faster loading
5. **Cache category tree** - Categories change infrequently, cache the tree
6. **Use slugs in URLs** - SEO-friendly URLs like `/categories/phones`

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Cannot delete category | Check for subcategories and products first |
| Category not showing | Verify `is_active` is true |
| Wrong product count | Clear cache and refresh |
| Slug conflict | System auto-generates unique slug with UUID |

---

## Implementation Summary

✅ **Features Implemented:**
- Hierarchical category structure
- Auto-generated slugs
- Product count per category
- Breadcrumb navigation
- Category tree API
- Featured categories
- Search functionality
- Admin CRUD operations
- Soft delete support

✅ **Files Created:**
- `internal/repository/category_repository.go`
- `internal/service/category_service.go`
- `internal/handler/category_handler.go`
- `api/routes_category.go`

✅ **Database:**
- Auto-migrated `categories` table
- Foreign key to `products` table
- Indexed for performance
