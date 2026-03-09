# Product Module - Business Flow & Implementation

## PART 1 — Product Business Flow

### Seller Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    SELLER PRODUCT MANAGEMENT                     │
└─────────────────────────────────────────────────────────────────┘

1. SELLER LOGS IN
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Seller  │────────>│  Server  │────────>│ Database │
   │          │  POST   │          │  Verify │          │
   │          │  /login │          │  Token  │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ JWT token with role="seller" generated
   ✓ Token stored in Authorization header

2. SELLER CREATES PRODUCT
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Seller  │────────>│ Product  │────────>│ Database │
   │          │  POST   │ Service  │  INSERT │          │
   │          │ /create │          │ Product │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Basic product info saved
   ✓ Product status = "draft" or "pending_review"

3. SELLER UPLOADS PRODUCT IMAGES
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Seller  │────────>│  Upload  │────────>│   CDN/   │
   │          │  POST   │ Service  │  Store  │  Storage │
   │          │ /images │          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Images uploaded to storage
   ✓ Image URLs saved to ProductImages table
   ✓ First image marked as primary

4. SELLER DEFINES PRODUCT VARIANTS
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Seller  │────────>│ Product  │────────>│ Database │
   │          │  POST   │ Variant  │  INSERT │          │
   │          │/variants│ Service  │ Variants│          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Variants created (e.g., Size: S/M/L, Color: Red/Blue)
   ✓ Each variant has unique SKU, price, stock

5. SELLER SETS PRICE AND INVENTORY
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Seller  │────────>│Inventory │────────>│ Database │
   │          │  PUT    │ Service  │  UPDATE │          │
   │          │ /stock  │          │  Stock  │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Stock levels set for each variant
   ✓ Inventory tracked in ProductInventory table

6. SELLER PUBLISHES PRODUCT
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Seller  │────────>│ Product  │────────>│ Database │
   │          │  PATCH  │ Service  │  UPDATE │          │
   │          │/publish │          │ Status  │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Product status changed to "active"
   ✓ Product visible to customers
   ✓ Product indexed in search engine
```

### Customer Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    CUSTOMER PRODUCT BROWSING                     │
└─────────────────────────────────────────────────────────────────┘

1. CUSTOMER SEARCHES PRODUCT
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│  Search  │────────>│Elastic/  │
   │          │  GET   │ Service  │  Query  │  SQL     │
   │          │ /search │          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Search by keyword
   ✓ Results ranked by relevance
   ✓ Filters applied (category, price, etc.)

2. CUSTOMER BROWSES CATEGORIES
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│ Category │────────>│ Database │
   │          │  GET    │ Service  │  Query  │          │
   │          │/category│          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Category tree loaded
   ✓ Products in category displayed
   ✓ Subcategories shown

3. CUSTOMER VIEWS PRODUCT DETAILS
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│ Product  │────────>│ Database │
   │          │  GET    │ Service  │  Query  │          │
   │          │   /:id  │          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Product details loaded
   ✓ Images displayed
   ✓ Variants shown
   ✓ Reviews loaded
   ✓ View count incremented

4. CUSTOMER SELECTS VARIANT
   ┌──────────┐         ┌──────────┐
   │ Customer │────────>│  Frontend │
   │          │  Select │           │
   │          │ Variant │           │
   └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Variant selected (Size: L, Color: Blue)
   ✓ Price updated based on variant
   ✓ Stock checked for variant

5. CUSTOMER ADDS PRODUCT TO CART
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│   Cart   │────────>│ Database │
   │          │  POST   │ Service  │  INSERT │          │
   │          │  /cart  │          │  Item   │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Stock reserved (15 minutes)
   ✓ Item added to cart
   ✓ Cart total updated
```

---

## PART 2 — Database Tables & Relationships

### Entity Relationship Diagram

```
┌─────────────────┐       ┌─────────────────┐
│   Categories    │       │      Shops      │
│─────────────────│       │─────────────────│
│ id (PK)         │       │ id (PK)         │
│ parent_id (FK)  │       │ user_id (FK)    │
│ name            │       │ name            │
│ slug            │       └────────┬────────┘
│ level           │                │
└────────┬────────┘                │
         │                         │
         │ 1:N                     │ 1:N
         ▼                         ▼
┌─────────────────────────────────────────┐
│              Products                    │
│─────────────────────────────────────────│
│ id (PK)                                  │
│ shop_id (FK) ────────────────────────────┤
│ category_id (FK) ────────────────────────┤
│ name                                     │
│ slug                                     │
│ description                              │
│ price                                    │
│ stock                                    │
│ status                                   │
└─────────────────┬───────────────────────┘
                  │
         ┌────────┼────────┐
         │        │        │
         │ 1:N    │ 1:N    │ 1:N
         ▼        ▼        ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ProductImages│ │ProductVar.  │ │ProductAttr. │
│─────────────│ │─────────────│ │─────────────│
│ id (PK)     │ │ id (PK)     │ │ id (PK)     │
│ product_id  │ │ product_id  │ │ product_id  │
│ url         │ │ sku         │ │ name        │
│ is_primary  │ │ price       │ │ value       │
│ sort_order  │ │ stock       │ │             │
└─────────────┘ │ attributes  │ └─────────────┘
                └─────────────┘
                        │
                        │ 1:N
                        ▼
                ┌─────────────┐
                │  Inventory  │
                │─────────────│
                │ id (PK)     │
                │ variant_id  │
                │ quantity    │
                │ reserved    │
                └─────────────┘
```

### Table Relationships

| Parent Table | Child Table | Relationship | Description |
|--------------|-------------|--------------|-------------|
| Categories | Products | 1:N | One category has many products |
| Categories | Categories | 1:N | Self-referential (parent-child) |
| Shops | Products | 1:N | One shop has many products |
| Products | ProductImages | 1:N | One product has many images |
| Products | ProductVariants | 1:N | One product has many variants |
| Products | ProductAttributes | 1:N | One product has many attributes |
| ProductVariants | Inventory | 1:1 | One variant has one inventory record |

---

## PART 3-13 — Implementation

The complete implementation follows in the code files:

### Files Created:

1. **Models** (`internal/domain/model/product_enhanced.go`)
   - Category, Product, ProductImage, ProductVariant
   - ProductAttribute, ProductAttributeValue, Inventory

2. **Repository** (`internal/repository/product_repository_enhanced.go`)
   - All CRUD operations
   - Search and filter methods
   - Inventory management

3. **Service** (`internal/service/product_service_enhanced.go`)
   - Business logic
   - Validation
   - Inventory operations

4. **Handler** (`internal/handler/product_handler_enhanced.go`)
   - REST API endpoints
   - Image upload
   - Search and filter

5. **Routes** (Updated in `api/routes_enhanced.go`)
   - Product routes with RBAC

---

## Key Features Implemented

### ✅ Product Management
- Create, update, delete products
- Product status workflow (draft → pending → active)
- Product visibility control

### ✅ Category Management
- Hierarchical categories (parent-child)
- Category-based product filtering
- SEO-friendly slugs

### ✅ Product Images
- Multiple images per product
- Primary image selection
- Image upload to server/CDN
- Image ordering

### ✅ Product Variants
- Multiple variant types (size, color, etc.)
- Variant-specific pricing
- Variant-specific stock
- Unique SKU per variant

### ✅ Inventory Management
- Stock tracking per variant
- Stock reservation on cart add
- Stock release on order cancel
- Low stock alerts

### ✅ Search & Filter
- Keyword search
- Category filter
- Price range filter
- Multiple sort options
- Pagination

### ✅ Security
- Seller can only manage own products
- Admin can manage all products
- Customer can only view products
- Role-based middleware

### ✅ Scalability
- Indexed queries
- Pagination support
- Caching ready
- Search engine ready

---

## API Endpoints Summary

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| POST | /api/products | Yes | Seller | Create product |
| PUT | /api/products/:id | Yes | Seller | Update product |
| DELETE | /api/products/:id | Yes | Admin | Delete product |
| GET | /api/products | No | Public | List products |
| GET | /api/products/:id | No | Public | Get product details |
| GET | /api/products/search | No | Public | Search products |
| GET | /api/products/category/:id | No | Public | Get by category |
| POST | /api/products/:id/images | Yes | Seller | Upload images |
| PUT | /api/products/:id/variants | Yes | Seller | Update variants |
| PUT | /api/products/:id/stock | Yes | Seller | Update stock |
| PATCH | /api/products/:id/publish | Yes | Seller | Publish product |

---

This product module is production-ready and can scale to millions of products with proper infrastructure.
