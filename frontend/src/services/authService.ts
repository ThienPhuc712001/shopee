import apiClient from './api';
import type { AuthResponse, LoginInput, RegisterInput, User } from '../types';

export const authService = {
  // Login
  async login(data: LoginInput): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/login', data);
    if (response.success) {
      // Store tokens
      localStorage.setItem('access_token', response.data.token);
      localStorage.setItem('refresh_token', response.data.refresh_token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  },

  // Register
  async register(data: RegisterInput): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/register', data);
    if (response.success) {
      localStorage.setItem('access_token', response.data.token);
      localStorage.setItem('refresh_token', response.data.refresh_token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  },

  // Logout
  async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout');
    } finally {
      this.clearAuth();
    }
  },

  // Refresh token
  async refreshToken(): Promise<AuthResponse> {
    const refreshToken = localStorage.getItem('refresh_token');
    const response = await apiClient.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken });
    if (response.success) {
      localStorage.setItem('access_token', response.data.token);
      localStorage.setItem('refresh_token', response.data.refresh_token);
    }
    return response.data;
  },

  // Get current user
  async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<User>('/auth/me');
    return response.data;
  },

  // Update profile
  async updateProfile(data: Partial<User>): Promise<User> {
    const response = await apiClient.put<User>('/auth/profile', data);
    if (response.success) {
      localStorage.setItem('user', JSON.stringify(response.data));
    }
    return response.data;
  },

  // Change password
  async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    await apiClient.put('/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword,
    });
  },

  // Clear auth data
  clearAuth(): void {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    apiClient.clearAuthTokens();
  },

  // Get stored user
  getStoredUser(): User | null {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch {
        return null;
      }
    }
    return null;
  },

  // Check if authenticated
  isAuthenticated(): boolean {
    return !!localStorage.getItem('access_token');
  },
};

export default authService;
