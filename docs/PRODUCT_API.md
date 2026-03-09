# Product Module API Documentation

## Base URL
```
http://localhost:8080/api
```

---

## Product Endpoints

### 1. Create Product
**POST** `/api/products`

Create a new product (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "Premium Wireless Headphones",
  "description": "High-quality wireless headphones with noise cancellation...",
  "short_description": "Premium wireless headphones with ANC",
  "price": 299.99,
  "original_price": 399.99,
  "stock": 100,
  "category_id": 1,
  "brand": "AudioTech",
  "weight": 250,
  "dimensions": "20x15x8",
  "warranty_period": "12 months",
  "return_days": 30,
  "tags": ["wireless", "headphones", "audio"],
  "meta_title": "Premium Wireless Headphones - AudioTech",
  "meta_description": "Buy premium wireless headphones with active noise cancellation"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "product": {
      "id": 1,
      "name": "Premium Wireless Headphones",
      "slug": "premium-wireless-headphones",
      "price": 299.99,
      "stock": 100,
      "status": "draft",
      "shop_id": 1,
      "category_id": 1,
      "created_at": "2024-01-15T10:30:00Z"
    }
  },
  "message": "Product created successfully"
}
```

---

### 2. Update Product
**PUT** `/api/products/:id`

Update an existing product (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "Updated Product Name",
  "price": 249.99,
  "stock": 150,
  "description": "Updated description..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "product": {
      "id": 1,
      "name": "Updated Product Name",
      "price": 249.99,
      ...
    }
  },
  "message": "Product updated successfully"
}
```

---

### 3. Delete Product
**DELETE** `/api/products/:id`

Delete a product (Admin or owner).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

---

### 4. Get Product by ID
**GET** `/api/products/:id`

Get product details by ID.

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "product": {
      "id": 1,
      "name": "Premium Wireless Headphones",
      "slug": "premium-wireless-headphones",
      "description": "High-quality wireless headphones...",
      "price": 299.99,
      "original_price": 399.99,
      "discount_percent": 25,
      "stock": 100,
      "available_stock": 95,
      "reserved_stock": 5,
      "sold_count": 250,
      "view_count": 1500,
      "rating_avg": 4.5,
      "rating_count": 45,
      "shop": {
        "id": 1,
        "name": "AudioTech Official Store",
        "rating": 4.8
      },
      "category": {
        "id": 1,
        "name": "Electronics",
        "slug": "electronics"
      },
      "images": [
        {
          "id": 1,
          "url": "/uploads/products/1_123456_image.jpg",
          "is_primary": true,
          "sort_order": 0
        }
      ],
      "variants": [
        {
          "id": 1,
          "name": "Black",
          "sku": "WH-001-BLK",
          "price": 299.99,
          "stock": 50
        }
      ]
    }
  }
}
```

---

### 5. Get All Products
**GET** `/api/products`

Get all products with pagination.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| limit | int | 20 | Items per page (max: 100) |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": 1,
        "name": "Premium Wireless Headphones",
        "price": 299.99,
        "rating_avg": 4.5,
        "sold_count": 250,
        "shop": {...},
        "images": [...]
      }
    ]
  },
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

### 6. Search Products
**GET** `/api/products/search`

Search products with filters.

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword |
| category_id | int | Filter by category |
| shop_id | int | Filter by shop |
| min_price | float | Minimum price |
| max_price | float | Maximum price |
| min_rating | float | Minimum rating |
| brands | string | Comma-separated brand names |
| sort_by | string | Sort field (price, rating, sold_count, created_at) |
| sort_order | string | Sort order (asc, desc) |
| page | int | Page number |
| limit | int | Items per page |

**Example Request:**
```
GET /api/products/search?keyword=headphones&min_price=100&max_price=500&sort_by=rating&sort_order=desc&page=1&limit=20
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [...]
  },
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

---

### 7. Get Products by Category
**GET** `/api/products/category/:id`

Get products in a specific category.

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [...]
  },
  "meta": {...}
}
```

---

### 8. Get Featured Products
**GET** `/api/products/featured`

Get featured products.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Number of products (max: 50) |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": 1,
        "name": "Featured Product",
        "is_featured": true,
        ...
      }
    ]
  }
}
```

---

### 9. Get Best Sellers
**GET** `/api/products/best-sellers`

Get best selling products.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Number of products (max: 50) |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": 1,
        "name": "Best Seller",
        "sold_count": 5000,
        ...
      }
    ]
  }
}
```

---

### 10. Upload Product Images
**POST** `/api/products/:id/images`

Upload images for a product (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

**Request Body (FormData):**
```
images: [file1, file2, file3]
```

**Image Requirements:**
- Max file size: 5MB
- Allowed formats: JPG, JPEG, PNG, GIF, WEBP
- Maximum images: 9 per product

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "images": [
      {
        "id": 1,
        "url": "/uploads/products/1_123456_image.jpg",
        "is_primary": true,
        "sort_order": 0
      }
    ]
  },
  "message": "Images uploaded successfully"
}
```

---

### 11. Create Product Variant
**POST** `/api/products/:id/variants`

Create a variant for a product (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "sku": "WH-001-BLK",
  "name": "Black",
  "price": 299.99,
  "original_price": 399.99,
  "stock": 50,
  "attributes": {
    "color": "Black",
    "size": "Standard"
  },
  "image_url": "/uploads/products/variant_black.jpg"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "variant": {
      "id": 1,
      "sku": "WH-001-BLK",
      "name": "Black",
      "price": 299.99,
      "stock": 50,
      "attributes": "{\"color\":\"Black\",\"size\":\"Standard\"}"
    }
  },
  "message": "Variant created successfully"
}
```

---

### 12. Get Product Variants
**GET** `/api/products/:id/variants`

Get all variants for a product.

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "variants": [
      {
        "id": 1,
        "sku": "WH-001-BLK",
        "name": "Black",
        "price": 299.99,
        "stock": 50,
        "attributes": "{\"color\":\"Black\"}"
      },
      {
        "id": 2,
        "sku": "WH-001-WHT",
        "name": "White",
        "price": 299.99,
        "stock": 30,
        "attributes": "{\"color\":\"White\"}"
      }
    ]
  }
}
```

---

### 13. Update Product Stock
**PUT** `/api/products/:id/stock`

Update stock for a product or variant (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "quantity": 100,
  "variant_id": 1
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Stock updated successfully"
}
```

---

### 14. Publish Product
**PATCH** `/api/products/:id/publish`

Publish a product (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Product published successfully"
}
```

---

### 15. Unpublish Product
**PATCH** `/api/products/:id/unpublish`

Unpublish a product (Seller only).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Product unpublished successfully"
}
```

---

## Category Endpoints

### 1. Get All Categories
**GET** `/api/categories`

Get all categories.

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "Electronics",
        "slug": "electronics",
        "level": 0,
        "parent_id": null,
        "children": [
          {
            "id": 2,
            "name": "Audio",
            "slug": "audio",
            "level": 1
          }
        ]
      }
    ]
  }
}
```

---

### 2. Get Category by ID
**GET** `/api/categories/:id`

Get category details.

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "category": {
      "id": 1,
      "name": "Electronics",
      "slug": "electronics",
      "description": "Electronic devices and accessories",
      "image_url": "/uploads/categories/electronics.jpg",
      "parent": null,
      "children": [...]
    }
  }
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error |

---

## Error Responses

```json
{
  "success": false,
  "error": "Error message here"
}
```

**Common Errors:**
- "Product not found"
- "Invalid product ID"
- "You can only update your own products"
- "Insufficient stock"
- "Invalid image format"
- "Image file too large"
- "Too many images (max 9)"
- "SKU already exists"

---

## cURL Examples

### Create Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Headphones",
    "price": 299.99,
    "stock": 100,
    "category_id": 1
  }'
```

### Search Products
```bash
curl -X GET "http://localhost:8080/api/products/search?keyword=headphones&min_price=100&max_price=500"
```

### Upload Images
```bash
curl -X POST http://localhost:8080/api/products/1/images \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.jpg"
```

### Update Stock
```bash
curl -X PUT http://localhost:8080/api/products/1/stock \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 150}'
```
