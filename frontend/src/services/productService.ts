import apiClient from './api';
import type { Product, ProductFilters, Category, Shop, Review, PaginatedResponse } from '../types';

export const productService = {
  // Get products with filters
  async getProducts(filters?: ProductFilters): Promise<PaginatedResponse<Product>> {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, String(value));
        }
      });
    }
    const response = await apiClient.get<PaginatedResponse<Product>>(`/products?${params.toString()}`);
    return response.data;
  },

  // Get featured products
  async getFeaturedProducts(limit: number = 10): Promise<Product[]> {
    const response = await apiClient.get<{ products: Product[] }>(`/products/featured?limit=${limit}`);
    return response.data.products;
  },

  // Get flash sale products
  async getFlashSaleProducts(): Promise<Product[]> {
    const response = await apiClient.get<{ products: Product[] }>('/products/flash-sale');
    return response.data.products;
  },

  // Get product by ID
  async getProductById(id: number): Promise<Product> {
    const response = await apiClient.get<Product>(`/products/${id}`);
    return response.data;
  },

  // Get product by slug
  async getProductBySlug(slug: string): Promise<Product> {
    const response = await apiClient.get<Product>(`/products/slug/${slug}`);
    return response.data;
  },

  // Get products by category
  async getProductsByCategory(categoryId: number, page: number = 1, limit: number = 20): Promise<PaginatedResponse<Product>> {
    const response = await apiClient.get<PaginatedResponse<Product>>(
      `/products?category_id=${categoryId}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Get products by shop
  async getProductsByShop(shopId: number, page: number = 1, limit: number = 20): Promise<PaginatedResponse<Product>> {
    const response = await apiClient.get<PaginatedResponse<Product>>(
      `/products?shop_id=${shopId}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Search products
  async searchProducts(query: string, page: number = 1, limit: number = 20): Promise<PaginatedResponse<Product>> {
    const response = await apiClient.get<PaginatedResponse<Product>>(
      `/products/search?q=${encodeURIComponent(query)}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Get categories
  async getCategories(): Promise<Category[]> {
    const response = await apiClient.get<{ categories: Category[] }>('/categories');
    return response.data.categories;
  },

  // Get category by ID
  async getCategoryById(id: number): Promise<Category> {
    const response = await apiClient.get<Category>(`/categories/${id}`);
    return response.data;
  },

  // Get category tree (for nested categories)
  async getCategoryTree(): Promise<Category[]> {
    const response = await apiClient.get<{ categories: Category[] }>('/categories/tree');
    return response.data.categories;
  },

  // Get product reviews
  async getProductReviews(productId: number, page: number = 1, limit: number = 10): Promise<PaginatedResponse<Review>> {
    const response = await apiClient.get<PaginatedResponse<Review>>(
      `/products/${productId}/reviews?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Create product review
  async createReview(productId: number, data: { rating: number; comment: string; images?: string[] }): Promise<Review> {
    const response = await apiClient.post<Review>(`/products/${productId}/reviews`, data);
    return response.data;
  },

  // Get shops
  async getShops(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Shop>> {
    const response = await apiClient.get<PaginatedResponse<Shop>>(`/shops?page=${page}&limit=${limit}`);
    return response.data;
  },

  // Get shop by ID
  async getShopById(id: number): Promise<Shop> {
    const response = await apiClient.get<Shop>(`/shops/${id}`);
    return response.data;
  },

  // Get related products
  async getRelatedProducts(productId: number, limit: number = 8): Promise<Product[]> {
    const response = await apiClient.get<{ products: Product[] }>(`/products/${productId}/related?limit=${limit}`);
    return response.data.products;
  },
};

export default productService;
