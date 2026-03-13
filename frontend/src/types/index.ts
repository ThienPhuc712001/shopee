// User types
export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name string;
  phone?: string;
  avatar?: string;
  role: 'customer' | 'seller' | 'admin';
  created_at: string;
}

export interface Address {
  id: number;
  user_id: number;
  name: string;
  phone: string;
  address: string;
  ward?: string;
  district?: string;
  city: string;
  is_default: boolean;
}

// Product types
export interface Product {
  id: number;
  shop_id: number;
  name: string;
  description: string;
  price: number;
  original_price?: number;
  stock: number;
  category_id: number;
  brand?: string;
  images: ProductImage[];
  rating?: number;
  review_count?: number;
  is_featured?: boolean;
  status: 'active' | 'inactive' | 'draft';
  created_at: string;
  updated_at: string;
}

export interface ProductImage {
  id: number;
  product_id: number;
  url: string;
  is_primary: boolean;
}

export interface Category {
  id: number;
  parent_id?: number;
  name: string;
  slug: string;
  description?: string;
  icon_url?: string;
  image_url?: string;
  level: number;
  sort_order: number;
  is_active: boolean;
  product_count?: number;
  children?: Category[];
}

// Cart types
export interface Cart {
  id: number;
  user_id: number;
  items: CartItem[];
  total_items: number;
  subtotal: number;
}

export interface CartItem {
  id: number;
  product_id: number;
  product: Product;
  quantity: number;
  subtotal: number;
}

// Order types
export interface Order {
  id: number;
  order_number: string;
  user_id: number;
  status: OrderStatus;
  items: OrderItem[];
  subtotal: number;
  shipping_fee: number;
  discount: number;
  total: number;
  shipping_address: Address;
  payment_method: string;
  created_at: string;
}

export interface OrderItem {
  id: number;
  product_id: number;
  product: Product;
  quantity: number;
  price: number;
  subtotal: number;
}

export type OrderStatus = 
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'
  | 'refunded';

// Shop types
export interface Shop {
  id: number;
  user_id: number;
  name: string;
  slug: string;
  description?: string;
  logo?: string;
  cover_image?: string;
  phone: string;
  email: string;
  address?: string;
  status: 'active' | 'pending' | 'suspended';
  rating?: number;
  rating_count?: number;
  follower_count?: number;
  product_count?: number;
  total_sales?: number;
  created_at: string;
}

// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: 'Bearer';
  expires_in: number;
  user: User;
}

// Notification types
export interface Notification {
  id: number;
  user_id: number;
  title: string;
  message: string;
  type: 'order' | 'promotion' | 'system' | 'shipping';
  is_read: boolean;
  created_at: string;
}

// Coupon types
export interface Coupon {
  id: number;
  code: string;
  discount_type: 'percentage' | 'fixed';
  discount_value: number;
  min_order_value?: number;
  max_discount?: number;
  valid_from: string;
  valid_until: string;
  usage_limit?: number;
  used_count?: number;
  is_active: boolean;
}
