import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosError, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios';
import { API_BASE_URL, STORAGE_KEYS } from '../constants';
import { type ApiResponse, type ApiError } from '../types';

class ApiClient {
  private instance: AxiosInstance;

  constructor() {
    this.instance = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    // Request interceptor - Add auth token
    this.instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const token = this.getAccessToken();
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error: AxiosError) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor - Handle errors
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => response,
      async (error: AxiosError<ApiError>) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

        // Handle 401 Unauthorized - Token expired
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const refreshToken = localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
            if (refreshToken) {
              const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
                refresh_token: refreshToken,
              });

              const { access_token, refresh_token } = response.data.data;
              this.setTokens(access_token, refresh_token);

              if (originalRequest.headers) {
                originalRequest.headers.Authorization = `Bearer ${access_token}`;
              }
              return this.instance(originalRequest);
            }
          } catch (refreshError) {
            // Refresh failed - logout user
            this.clearAuth();
            window.location.href = '/login';
            return Promise.reject(refreshError);
          }
        }

        // Handle other errors
        const errorMessage = error.response?.data?.message || 'An error occurred';
        return Promise.reject({
          success: false,
          message: errorMessage,
          error: error.message,
          status: error.response?.status || 500,
        });
      }
    );
  }

  private getAccessToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  }

  private setTokens(accessToken: string, refreshToken: string): void {
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken);
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
  }

  private clearAuth(): void {
    localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER);
  }

  // HTTP Methods
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.instance.get<ApiResponse<T>>(url, config);
    return response.data;
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.instance.post<ApiResponse<T>>(url, data, config);
    return response.data;
  }

  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.instance.put<ApiResponse<T>>(url, data, config);
    return response.data;
  }

  async patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.instance.patch<ApiResponse<T>>(url, data, config);
    return response.data;
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.instance.delete<ApiResponse<T>>(url, config);
    return response.data;
  }

  // Auth methods
  setAuthTokens(accessToken: string, refreshToken: string): void {
    this.setTokens(accessToken, refreshToken);
  }

  clearAuthTokens(): void {
    this.clearAuth();
  }

  isAuthenticated(): boolean {
    return !!this.getAccessToken();
  }
}

export const apiClient = new ApiClient();
export default apiClient;
