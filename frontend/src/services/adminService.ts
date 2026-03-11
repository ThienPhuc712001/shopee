import apiClient from './api';
import type { User, Product, Order, AdminStats, AdminUser, Category, Shop, PaginatedResponse } from '../types';

export interface AdminProductInput {
  name: string;
  category_id: number;
  description: string;
  short_description?: string;
  price: number;
  original_price?: number;
  stock: number;
  sku?: string;
  brand?: string;
  images?: { url: string; is_primary: boolean }[];
  status?: string;
}

export const adminService = {
  // Dashboard
  async getDashboardStats(): Promise<AdminStats> {
    const response = await apiClient.get<AdminStats>('/admin/dashboard');
    return response.data;
  },

  // Analytics
  async getAnalytics(period: string = '30d'): Promise<{
    sales: { date: string; value: number }[];
    orders: { date: string; value: number }[];
    users: { date: string; value: number }[];
  }> {
    const response = await apiClient.get(`/admin/analytics?period=${period}`);
    return response.data as { sales: { date: string; value: number }[]; orders: { date: string; value: number }[]; users: { date: string; value: number }[] };
  },

  // Users Management
  async getUsers(page: number = 1, limit: number = 20): Promise<PaginatedResponse<User>> {
    const response = await apiClient.get<PaginatedResponse<User>>(`/admin/users?page=${page}&limit=${limit}`);
    return response.data;
  },

  async getUserById(id: number): Promise<User> {
    const response = await apiClient.get<User>(`/admin/users/${id}`);
    return response.data;
  },

  async banUser(id: number, reason: string): Promise<User> {
    const response = await apiClient.post<User>(`/admin/users/${id}/ban`, { reason });
    return response.data;
  },

  async unbanUser(id: number): Promise<User> {
    const response = await apiClient.post<User>(`/admin/users/${id}/unban`);
    return response.data;
  },

  async deleteUser(id: number): Promise<void> {
    await apiClient.delete(`/admin/users/${id}`);
  },

  // Admin Users Management
  async getAdminUsers(): Promise<AdminUser[]> {
    const response = await apiClient.get<{ users: AdminUser[] }>('/admin/admin-users');
    return response.data.users;
  },

  async createAdminUser(data: { email: string; password: string; role_id: number; first_name: string; last_name: string }): Promise<AdminUser> {
    const response = await apiClient.post<AdminUser>('/admin/admin-users', data);
    return response.data;
  },

  async updateAdminUser(id: number, data: Partial<AdminUser>): Promise<AdminUser> {
    const response = await apiClient.put<AdminUser>(`/admin/admin-users/${id}`, data);
    return response.data;
  },

  async deleteAdminUser(id: number): Promise<void> {
    await apiClient.delete(`/admin/admin-users/${id}`);
  },

  // Products Management
  async getProducts(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Product>> {
    const response = await apiClient.get<PaginatedResponse<Product>>(`/admin/products?page=${page}&limit=${limit}`);
    return response.data;
  },

  async getProductById(id: number): Promise<Product> {
    const response = await apiClient.get<Product>(`/admin/products/${id}`);
    return response.data;
  },

  async createProduct(data: AdminProductInput): Promise<Product> {
    const response = await apiClient.post<Product>('/admin/products', data);
    return response.data;
  },

  async updateProduct(id: number, data: Partial<AdminProductInput>): Promise<Product> {
    const response = await apiClient.put<Product>(`/admin/products/${id}`, data);
    return response.data;
  },

  async deleteProduct(id: number): Promise<void> {
    await apiClient.delete(`/admin/products/${id}`);
  },

  async updateProductStatus(id: number, status: string): Promise<Product> {
    const response = await apiClient.patch<Product>(`/admin/products/${id}/status`, { status });
    return response.data;
  },

  // Orders Management
  async getOrders(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Order>> {
    const response = await apiClient.get<PaginatedResponse<Order>>(`/admin/orders?page=${page}&limit=${limit}`);
    return response.data;
  },

  async getOrderById(id: number): Promise<Order> {
    const response = await apiClient.get<Order>(`/admin/orders/${id}`);
    return response.data;
  },

  async updateOrderStatus(id: number, status: string): Promise<Order> {
    const response = await apiClient.patch<Order>(`/admin/orders/${id}/status`, { status });
    return response.data;
  },

  async cancelOrder(id: number, reason: string): Promise<Order> {
    const response = await apiClient.post<Order>(`/admin/orders/${id}/cancel`, { reason });
    return response.data;
  },

  async refundOrder(id: number, amount: number, reason: string): Promise<Order> {
    const response = await apiClient.post<Order>(`/admin/orders/${id}/refund`, { amount, reason });
    return response.data;
  },

  // Shops Management
  async getShops(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Shop>> {
    const response = await apiClient.get<PaginatedResponse<Shop>>(`/admin/shops?page=${page}&limit=${limit}`);
    return response.data;
  },

  async approveShop(id: number): Promise<Shop> {
    const response = await apiClient.post<Shop>(`/admin/shops/${id}/approve`);
    return response.data;
  },

  async rejectShop(id: number, reason: string): Promise<Shop> {
    const response = await apiClient.post<Shop>(`/admin/shops/${id}/reject`, { reason });
    return response.data;
  },

  async suspendShop(id: number, reason: string): Promise<Shop> {
    const response = await apiClient.post<Shop>(`/admin/shops/${id}/suspend`, { reason });
    return response.data;
  },

  // Categories Management
  async getCategories(): Promise<Category[]> {
    const response = await apiClient.get<{ categories: Category[] }>('/admin/categories');
    return response.data.categories;
  },

  async createCategory(data: { name: string; slug: string; parent_id?: number; description?: string }): Promise<Category> {
    const response = await apiClient.post<Category>('/admin/categories', data);
    return response.data;
  },

  async updateCategory(id: number, data: Partial<Category>): Promise<Category> {
    const response = await apiClient.put<Category>(`/admin/categories/${id}`, data);
    return response.data;
  },

  async deleteCategory(id: number): Promise<void> {
    await apiClient.delete(`/admin/categories/${id}`);
  },

  // System Settings
  async getSettings(): Promise<Record<string, string>> {
    const response = await apiClient.get<Record<string, string>>('/admin/settings');
    return response.data;
  },

  async updateSettings(settings: Record<string, string>): Promise<void> {
    await apiClient.put('/admin/settings', settings);
  },
};

export default adminService;
