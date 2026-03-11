# E-Commerce Frontend

A modern, responsive React/TypeScript frontend for a large-scale e-commerce platform similar to Shopee or Lazada.

## Tech Stack

- **Framework:** React 18
- **Language:** TypeScript
- **Routing:** React Router v6
- **State Management:** Redux Toolkit
- **UI Framework:** Tailwind CSS
- **HTTP Client:** Axios
- **Form Handling:** React Hook Form + Zod validation
- **Icons:** Heroicons

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── common/         # Generic components (Button, Input, Modal, etc.)
│   ├── layout/         # Layout components (Navbar, Footer)
│   ├── product/        # Product-related components
│   ├── cart/           # Cart-related components
│   ├── checkout/       # Checkout components
│   ├── order/          # Order components
│   └── admin/          # Admin components
├── pages/              # Page components
│   ├── auth/           # Login, Register pages
│   ├── product/        # Products list, Product detail
│   ├── cart/           # Cart page
│   ├── checkout/       # Checkout pages
│   ├── user/           # User dashboard pages
│   └── admin/          # Admin dashboard pages
├── services/           # API service layer
│   ├── api.ts          # Axios instance & interceptors
│   ├── authService.ts  # Auth API calls
│   ├── productService.ts
│   ├── cartService.ts
│   ├── orderService.ts
│   ├── userService.ts
│   └── adminService.ts
├── store/              # Redux store
│   └── slices/         # Redux slices
│       ├── authSlice.ts
│       ├── cartSlice.ts
│       ├── productsSlice.ts
│       └── ordersSlice.ts
├── hooks/              # Custom React hooks
├── types/              # TypeScript type definitions
├── constants/          # App constants (routes, config)
├── utils/              # Utility functions
└── assets/             # Static assets
```

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Create .env file (already created)
# VITE_API_URL=http://localhost:8080/api

# Start development server
npm run dev
```

### Build

```bash
# Production build
npm run build

# Preview production build
npm run preview
```

## Features

### PART 1 - Project Architecture
- ✅ Modular folder structure
- ✅ TypeScript for type safety
- ✅ Tailwind CSS for styling
- ✅ Component-based architecture

### PART 2 - Global State Management
- ✅ Redux Toolkit for state management
- ✅ Auth slice (user, token, authentication state)
- ✅ Products slice (products, categories, filters)
- ✅ Cart slice (cart items, totals)
- ✅ Orders slice (orders, checkout state)

### PART 3 - API Integration
- ✅ Axios HTTP client
- ✅ Request/Response interceptors
- ✅ JWT token handling
- ✅ Automatic token refresh
- ✅ Error handling

### PART 4 - Authentication
- ✅ Login page with form validation
- ✅ Register page with form validation
- ✅ JWT token storage (localStorage)
- ✅ Protected routes
- ✅ Auto-redirect after login

### PART 5 - Home Page
- ✅ Hero section
- ✅ Category browsing
- ✅ Featured products
- ✅ Flash sale section
- ✅ Search functionality

### PART 6 - Product List Page
- ✅ Product grid layout
- ✅ Category filtering
- ✅ Price range filtering
- ✅ Rating filtering
- ✅ Sorting options
- ✅ Pagination

### PART 7 - Product Detail Page
- ✅ Image gallery
- ✅ Product information
- ✅ Variant selection
- ✅ Quantity selector
- ✅ Add to cart
- ✅ Related products

### PART 8 - Cart Page
- ✅ Cart items display
- ✅ Quantity update
- ✅ Item removal
- ✅ Price calculation
- ✅ Coupon code input
- ✅ Checkout redirect

### PART 9 - Checkout Page
- ✅ Multi-step checkout
- ✅ Shipping information form
- ✅ Payment method selection
- ✅ Order summary
- ✅ Order placement

### PART 10 - User Dashboard
- ✅ User profile page
- ✅ Order history
- ✅ Order tracking
- ✅ Address management (placeholder)
- ✅ Review history (placeholder)

### PART 11 - Admin Dashboard
- ✅ Admin layout structure
- ⏳ User management
- ⏳ Product management
- ⏳ Order management
- ⏳ Analytics charts

### PART 12 - UI Components
- ✅ Navbar
- ✅ Footer
- ✅ ProductCard
- ✅ CartItem
- ✅ OrderCard
- ✅ Pagination
- ✅ Modal
- ✅ LoadingSpinner
- ✅ Button
- ✅ Input

### PART 13 - Security
- ✅ Protected routes
- ✅ JWT token storage
- ✅ Token refresh mechanism
- ✅ Logout handling
- ✅ Session management

### PART 14 - Performance Optimization
- ✅ Lazy loading with React.lazy
- ✅ Code splitting by route
- ✅ Image lazy loading
- ✅ Memoization ready structure

### PART 15 - Example Implementation
- ✅ Product list page
- ✅ Product detail page
- ✅ Add to cart functionality
- ✅ Checkout flow

## API Integration

The frontend expects the following API endpoints:

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh token
- `GET /api/auth/me` - Get current user

### Products
- `GET /api/products` - List products (with filters)
- `GET /api/products/:id` - Get product by ID
- `GET /api/products/featured` - Featured products
- `GET /api/products/flash-sale` - Flash sale products
- `GET /api/categories` - List categories

### Cart
- `GET /api/cart` - Get user's cart
- `POST /api/cart/items` - Add item to cart
- `PUT /api/cart/items/:id` - Update cart item
- `DELETE /api/cart/items/:id` - Remove item

### Orders
- `POST /api/orders/checkout` - Create order from cart
- `GET /api/orders` - Get user's orders
- `GET /api/orders/:id` - Get order details
- `POST /api/orders/:id/cancel` - Cancel order

### User
- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update profile
- `GET /api/user/addresses` - Get addresses
- `POST /api/user/addresses` - Add address

## State Management

### Auth Slice
```typescript
{
  user: User | null,
  token: string | null,
  isAuthenticated: boolean,
  isLoading: boolean,
  error: string | null
}
```

### Products Slice
```typescript
{
  products: Product[],
  featuredProducts: Product[],
  flashSaleProducts: Product[],
  categories: Category[],
  pagination: {...},
  filters: ProductFilters,
  isLoading: boolean,
  error: string | null
}
```

### Cart Slice
```typescript
{
  cart: Cart | null,
  items: CartItem[],
  itemCount: number,
  isLoading: boolean,
  error: string | null
}
```

### Orders Slice
```typescript
{
  orders: Order[],
  currentOrder: Order | null,
  stats: OrderStats | null,
  isCheckingOut: boolean,
  isLoading: boolean,
  error: string | null
}
```

## Environment Variables

```env
VITE_API_URL=http://localhost:8080/api
```

## Styling

The project uses Tailwind CSS with custom configuration:

- Primary color: Red (#ef4444)
- Secondary color: Slate (#64748b)
- Custom button variants
- Custom input styles
- Card components

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

MIT
