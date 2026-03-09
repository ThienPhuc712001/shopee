# E-Commerce Platform Backend

A production-ready, scalable e-commerce backend built with Golang, Gin, GORM, and SQL Server following Clean Architecture principles.

## Project Structure

```
cmd/
  server/
    main.go              # Application entry point

internal/
  domain/
    model/               # Domain entities/models
    repository/          # Repository interfaces

  service/               # Business logic layer
  handler/               # HTTP handlers

pkg/
  config/                # Configuration management
  database/              # Database connection
  middleware/            # HTTP middleware
  utils/                 # Utility functions
  response/              # Standard API responses

api/
  auth/                  # Authentication routes
  users/                 # User management routes
  shops/                 # Shop management routes
  products/              # Product management routes
  cart/                  # Shopping cart routes
  orders/                # Order management routes
  payments/              # Payment processing routes
  reviews/               # Review management routes
  admin/                 # Admin panel routes
```

## Folder Purposes

| Folder | Purpose |
|--------|---------|
| `cmd/server` | Application entry point, server initialization |
| `internal/domain/model` | Domain entities and data structures |
| `internal/domain/repository` | Repository interfaces (contracts) |
| `internal/service` | Business logic and use cases |
| `internal/handler` | HTTP request handlers |
| `pkg/config` | Configuration loading and management |
| `pkg/database` | Database connection and migration |
| `pkg/middleware` | HTTP middleware (auth, logging, CORS) |
| `pkg/utils` | Helper functions and utilities |
| `pkg/response` | Standardized API response formats |
| `api/*` | Route definitions for each module |

## Tech Stack

- **Language**: Golang
- **Framework**: Gin
- **Database**: Microsoft SQL Server
- **ORM**: GORM
- **Authentication**: JWT
- **Architecture**: Clean Architecture

## Getting Started

### Prerequisites

- Go 1.21+
- SQL Server 2019+
- Git

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd TMDT

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run the application
go run cmd/server/main.go
```

### Configuration

Create a `.env` file with the following variables:

```env
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=1433
DB_NAME=ecommerce
DB_USER=sa
DB_PASSWORD=your_password

JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

BCRYPT_COST=10
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login

### Products
- `GET /api/products` - List products
- `GET /api/products/:id` - Get product details
- `POST /api/products` - Create product (authenticated)

### Cart
- `POST /api/cart/add` - Add item to cart
- `GET /api/cart` - Get cart items
- `DELETE /api/cart/:id` - Remove item from cart

### Orders
- `POST /api/orders` - Create order
- `GET /api/orders/:id` - Get order details
- `GET /api/orders` - List user orders

## License

MIT
