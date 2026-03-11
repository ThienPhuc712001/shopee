import apiClient from './api';
import type { User, Address, Review, PaginatedResponse } from '../types';

export const userService = {
  // Get user profile
  async getProfile(): Promise<User> {
    const response = await apiClient.get<User>('/user/profile');
    return response.data;
  },

  // Update user profile
  async updateProfile(data: Partial<User>): Promise<User> {
    const response = await apiClient.put<User>('/user/profile', data);
    return response.data;
  },

  // Get user addresses
  async getAddresses(): Promise<Address[]> {
    const response = await apiClient.get<{ addresses: Address[] }>('/user/addresses');
    return response.data.addresses;
  },

  // Add address
  async addAddress(data: Omit<Address, 'id' | 'user_id'>): Promise<Address> {
    const response = await apiClient.post<Address>('/user/addresses', data);
    return response.data;
  },

  // Update address
  async updateAddress(id: number, data: Partial<Address>): Promise<Address> {
    const response = await apiClient.put<Address>(`/user/addresses/${id}`, data);
    return response.data;
  },

  // Delete address
  async deleteAddress(id: number): Promise<void> {
    await apiClient.delete(`/user/addresses/${id}`);
  },

  // Set default address
  async setDefaultAddress(id: number): Promise<Address[]> {
    const response = await apiClient.post<{ addresses: Address[] }>(`/user/addresses/${id}/default`);
    return response.data.addresses;
  },

  // Get user reviews
  async getReviews(page: number = 1, limit: number = 10): Promise<PaginatedResponse<Review>> {
    const response = await apiClient.get<PaginatedResponse<Review>>(
      `/user/reviews?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Get user's reviewed products
  async getReviewedProducts(): Promise<{ product_id: number; reviewed: boolean }[]> {
    const response = await apiClient.get('/user/reviewed-products');
    return response.data as { product_id: number; reviewed: boolean }[];
  },

  // Delete user account
  async deleteAccount(): Promise<void> {
    await apiClient.delete('/user/account');
  },

  // Get user statistics
  async getStats(): Promise<{
    total_orders: number;
    total_spent: number;
    pending_orders: number;
    completed_orders: number;
    total_reviews: number;
  }> {
    const response = await apiClient.get('/user/stats');
    return response.data as { total_orders: number; total_spent: number; pending_orders: number; completed_orders: number; total_reviews: number };
  },
};

export default userService;
