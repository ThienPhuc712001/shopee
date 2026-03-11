export const API_BASE_URL = '/api';

export const ROUTES = {
  // Public routes
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  PRODUCTS: '/products',
  PRODUCT_DETAIL: '/products/:id',
  SHOP: '/shops/:id',
  
  // Cart & Checkout
  CART: '/cart',
  CHECKOUT: '/checkout',
  CHECKOUT_SUCCESS: '/checkout/success',
  
  // User routes
  USER_DASHBOARD: '/user',
  USER_PROFILE: '/user/profile',
  USER_ORDERS: '/user/orders',
  USER_ORDER_DETAIL: '/user/orders/:id',
  USER_REVIEWS: '/user/reviews',
  USER_ADDRESSES: '/user/addresses',
  USER_SETTINGS: '/user/settings',
  
  // Admin routes
  ADMIN_DASHBOARD: '/admin',
  ADMIN_USERS: '/admin/users',
  ADMIN_PRODUCTS: '/admin/products',
  ADMIN_ORDERS: '/admin/orders',
  ADMIN_CATEGORIES: '/admin/categories',
  ADMIN_SHOPS: '/admin/shops',
  ADMIN_ANALYTICS: '/admin/analytics',
  ADMIN_SETTINGS: '/admin/settings',
} as const;

export const USER_ROLES = {
  CUSTOMER: 'customer',
  SELLER: 'seller',
  ADMIN: 'admin',
} as const;

export const ORDER_STATUS = {
  PENDING: 'pending',
  PAID: 'paid',
  PROCESSING: 'processing',
  SHIPPED: 'shipped',
  DELIVERED: 'delivered',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const;

export const PAYMENT_STATUS = {
  PENDING: 'pending',
  PROCESSING: 'processing',
  PAID: 'paid',
  FAILED: 'failed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const;

export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'access_token',
  REFRESH_TOKEN: 'refresh_token',
  USER: 'user',
  CART: 'cart',
} as const;

export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_LIMIT: 20,
  LIMIT_OPTIONS: [10, 20, 50, 100],
} as const;

export const SORT_OPTIONS = [
  { value: 'newest', label: 'Newest' },
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
  { value: 'best_selling', label: 'Best Selling' },
  { value: 'top_rated', label: 'Top Rated' },
] as const;
