# E-Commerce Frontend - React + TypeScript + TailwindCSS

Modern, professional e-commerce web interface built with React, TypeScript, and TailwindCSS.

## 🚀 Tech Stack

- **Framework:** React 18
- **Language:** TypeScript
- **Styling:** TailwindCSS
- **State Management:** Redux Toolkit
- **Routing:** React Router v6
- **Icons:** Lucide React
- **Animations:** Framer Motion
- **HTTP Client:** Axios

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/          # Reusable components
│   │   ├── common/         # Common components (Header, Footer)
│   │   ├── product/        # Product-related components
│   │   ├── cart/           # Cart components
│   │   ├── checkout/       # Checkout components
│   │   └── admin/          # Admin components
│   ├── pages/              # Page components
│   │   ├── home/           # Homepage
│   │   ├── product/        # Product pages
│   │   ├── cart/           # Cart page
│   │   ├── checkout/       # Checkout page
│   │   ├── user/           # User dashboard
│   │   └── admin/          # Admin panel
│   ├── layouts/            # Layout components
│   ├── store/              # Redux store
│   │   └── slices/         # Redux slices
│   ├── services/           # API services
│   ├── types/              # TypeScript types
│   ├── utils/              # Utility functions
│   └── hooks/              # Custom hooks
├── public/                 # Static assets
└── package.json
```

## 🎨 Design System

### Colors
- **Primary:** `#FF5722` (Orange - for CTAs, highlights)
- **Secondary:** `#1E293B` (Dark Slate - for text, headers)
- **Background:** `#F8FAFC` (Light Gray - for page background)
- **Surface:** `#FFFFFF` (White - for cards)
- **Text:** `#111827` (Almost Black - for primary text)

### Typography
- **Font Family:** Inter
- **Font Weights:** 400 (Normal), 500 (Medium), 600 (Semibold), 700 (Bold)

### Spacing
Using 8px grid system: 4px, 8px, 12px, 16px, 24px, 32px, 48px, 64px

## 🛠️ Setup & Installation

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Install Dependencies

```bash
cd frontend
npm install
```

### Start Development Server

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## 🌐 API Integration

The frontend connects to the backend API at `http://localhost:8080/api`.

To change the API URL, create a `.env` file:

```env
VITE_API_URL=http://localhost:8080/api
```

## 📱 Responsive Breakpoints

- **Mobile:** < 640px
- **Tablet:** 640px - 1024px
- **Desktop:** > 1024px

## ✨ Features Implemented

### Homepage
- ✅ Hero Banner with auto-sliding carousel
- ✅ Categories Section with icons
- ✅ Flash Sale Section with countdown timer
- ✅ Featured Products Grid
- ✅ Recommended Products

### Product Pages
- ✅ Product Listing with filters
- ✅ Product Grid (responsive)
- ✅ Product Card with hover effects
- ✅ Product Detail Page (placeholder)

### Cart & Checkout
- ✅ Shopping Cart Page
- ✅ Checkout Page (placeholder)
- ✅ Order Summary

### User Dashboard
- ✅ Account Dashboard
- ✅ Orders Management
- ✅ Wishlist
- ✅ Addresses

### Admin Panel
- ✅ Admin Dashboard with stats
- ✅ Product Management
- ✅ Order Management
- ✅ User Management
- ✅ Analytics

## 🎯 Key Components

### Header
- Sticky navigation
- Search bar with category dropdown
- Cart icon with item count badge
- User menu with dropdown
- Mobile responsive menu

### Footer
- Company information
- Quick links
- Customer service links
- Social media icons
- Contact information

### Product Card
- Product image with hover zoom
- Discount badge
- Quick action buttons (wishlist, view)
- Add to cart button on hover
- Rating stars
- Price with discount
- Stock status

## 🔧 State Management

Redux slices:
- `authSlice` - User authentication state
- `productSlice` - Products and categories
- `cartSlice` - Shopping cart
- `orderSlice` - Orders

## 📊 API Endpoints Used

```typescript
// Auth
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/me

// Products
GET  /api/products
GET  /api/products/:id
GET  /api/products/featured
GET  /api/products/best-sellers
GET  /api/products/search
GET  /api/products/category/:id

// Categories
GET  /api/categories
GET  /api/categories/tree

// Cart
GET  /api/cart
POST /api/cart/add
PUT  /api/cart/items/:id
DELETE /api/cart/items/:id

// Orders
POST /api/orders
GET  /api/orders
GET  /api/orders/:id
```

## 🎨 Animation Guidelines

- Page transitions: Fade in
- Hover effects: Scale + shadow
- Cart animation: Slide up
- Button clicks: Scale down

## 📝 Code Style

- Use TypeScript for type safety
- Follow ESLint rules
- Use functional components with hooks
- Implement proper error handling
- Add loading states

## 🚀 Performance Optimizations

- Lazy loading for images
- Code splitting for routes
- Debounced search
- Memoized components
- Virtual scrolling for long lists

## 📱 Mobile Considerations

- Touch-friendly buttons (min 44px)
- Swipe gestures for carousels
- Collapsible filters
- Bottom navigation option
- Image optimization

## 🔐 Security

- JWT token storage in localStorage
- Automatic token refresh
- Protected routes
- Input validation
- XSS protection

## 📄 License

MIT License

---

**Built with ❤️ using React, TypeScript, and TailwindCSS**
