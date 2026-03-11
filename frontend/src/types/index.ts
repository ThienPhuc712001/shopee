// User types
export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
  avatar: string;
  role: 'customer' | 'seller' | 'admin';
  status: 'active' | 'inactive' | 'banned' | 'locked';
  email_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface AdminUser {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
  avatar_url: string;
  role_id: number;
  status: 'active' | 'inactive' | 'suspended';
  created_at: string;
  updated_at: string;
}

// Product types
export interface Category {
  id: number;
  parent_id: number | null;
  name: string;
  slug: string;
  description: string;
  icon_url: string;
  image_url: string;
  level: number;
  sort_order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  shop_id: number;
  category_id: number;
  name: string;
  slug: string;
  description: string;
  short_description: string;
  sku: string;
  brand: string;
  price: number;
  original_price: number;
  discount_percent: number;
  stock: number;
  reserved_stock: number;
  available_stock: number;
  sold_count: number;
  view_count: number;
  rating_avg: number;
  rating_count: number;
  status: 'draft' | 'pending' | 'active' | 'inactive' | 'banned' | 'out_of_stock';
  is_featured: boolean;
  is_flash_sale: boolean;
  flash_sale_price: number;
  flash_sale_start: string | null;
  flash_sale_end: string | null;
  images: ProductImage[];
  variants: ProductVariant[];
  category?: Category;
  shop?: Shop;
  created_at: string;
  updated_at: string;
}

export interface ProductImage {
  id: number;
  product_id: number;
  url: string;
  alt_text: string;
  is_primary: boolean;
  sort_order: number;
  created_at: string;
}

export interface ProductVariant {
  id: number;
  product_id: number;
  sku: string;
  name: string;
  price: number;
  original_price: number;
  stock: number;
  attributes: Record<string, string>;
  image_url: string;
  sort_order: number;
}

export interface Shop {
  id: number;
  user_id: number;
  name: string;
  slug: string;
  description: string;
  logo: string;
  cover_image: string;
  phone: string;
  email: string;
  address: string;
  status: 'pending' | 'active' | 'inactive' | 'suspended';
  verification_status: string;
  rating: number;
  rating_count: number;
  follower_count: number;
  product_count: number;
  total_sales: number;
  total_revenue: number;
  created_at: string;
  updated_at: string;
}

// Cart types
export interface Cart {
  id: number;
  user_id: number;
  total_items: number;
  subtotal: number;
  discount: number;
  total: number;
  currency: string;
  items: CartItem[];
  created_at: string;
  updated_at: string;
}

export interface CartItem {
  id: number;
  cart_id: number;
  product_id: number;
  variant_id: number | null;
  quantity: number;
  price: number;
  original_price: number;
  discount: number;
  subtotal: number;
  product_name: string;
  product_image: string;
  shop_id: number;
  product?: Product;
  variant?: ProductVariant;
  shop?: Shop;
  created_at: string;
  updated_at: string;
}

// Order types
export interface Order {
  id: number;
  order_number: string;
  user_id: number;
  shop_id: number;
  status: OrderStatus;
  payment_status: PaymentStatus;
  fulfillment_status: FulfillmentStatus;
  subtotal: number;
  shipping_fee: number;
  shipping_discount: number;
  product_discount: number;
  voucher_discount: number;
  tax_amount: number;
  total_amount: number;
  paid_amount: number;
  shipping_name: string;
  shipping_phone: string;
  shipping_address: string;
  shipping_ward: string;
  shipping_district: string;
  shipping_city: string;
  shipping_country: string;
  shipping_method: string;
  tracking_number: string;
  buyer_note: string;
  cancel_reason: string;
  items: OrderItem[];
  paid_at: string | null;
  shipped_at: string | null;
  delivered_at: string | null;
  cancelled_at: string | null;
  completed_at: string | null;
  created_at: string;
  updated_at: string;
}

export type OrderStatus = 'pending' | 'paid' | 'processing' | 'shipped' | 'delivered' | 'cancelled' | 'refunded';
export type PaymentStatus = 'pending' | 'processing' | 'paid' | 'failed' | 'cancelled' | 'refunded';
export type FulfillmentStatus = 'unfulfilled' | 'processing' | 'packed' | 'shipped' | 'delivered';

export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  variant_id: number | null;
  product_name: string;
  product_image: string;
  product_sku: string;
  quantity: number;
  price: number;
  original_price: number;
  discount: number;
  subtotal: number;
  tax_amount: number;
  final_amount: number;
  shop_id: number;
  fulfillment_status: FulfillmentStatus;
  product?: Product;
}

export interface Address {
  id: number;
  user_id: number;
  name: string;
  phone: string;
  street: string;
  ward: string;
  district: string;
  city: string;
  country: string;
  is_default: boolean;
}

// Review types
export interface Review {
  id: number;
  user_id: number;
  product_id: number;
  shop_id: number;
  order_id: number;
  rating: number;
  comment: string;
  images: string[];
  is_approved: boolean;
  helpful_count: number;
  user?: User;
  created_at: string;
  updated_at: string;
}

// Auth types
export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone: string;
}

export interface AuthResponse {
  user: User;
  token: string;
  refresh_token: string;
}

// Re-export API types
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ApiError {
  success: false;
  message: string;
  error: string;
  status: number;
}

// Filter types
export interface ProductFilters {
  category_id?: number;
  shop_id?: number;
  min_price?: number;
  max_price?: number;
  rating?: number;
  sort?: 'newest' | 'price_asc' | 'price_desc' | 'best_selling' | 'top_rated';
  search?: string;
  page?: number;
  limit?: number;
}

// Admin types
export interface AdminStats {
  total_users: number;
  total_sellers: number;
  total_products: number;
  total_orders: number;
  total_revenue: number;
  pending_orders: number;
  pending_refunds: number;
  active_users_24h: number;
  new_users_today: number;
}

export interface OrderStats {
  total_orders: number;
  total_revenue: number;
  pending_orders: number;
  completed_orders: number;
  cancelled_orders: number;
  average_order_value: number;
}
