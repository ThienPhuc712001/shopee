import apiClient from './api'
import { ApiResponse, Product, Category, PaginatedResponse } from '../types'

interface ProductFilters {
  category_id?: number
  keyword?: string
  min_price?: number
  max_price?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  page?: number
  limit?: number
}

export const productService = {
  getProducts: async (filters?: ProductFilters): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<Product>>>('/products', { params: filters })
    return response.data.data!
  },

  getProductById: async (id: number): Promise<Product> => {
    const response = await apiClient.get<ApiResponse<Product>>(`/products/${id}`)
    return response.data.data!
  },

  getFeaturedProducts: async (): Promise<Product[]> => {
    const response = await apiClient.get<ApiResponse<Product[]>>('/products/featured')
    return response.data.data!
  },

  getBestSellers: async (): Promise<Product[]> => {
    const response = await apiClient.get<ApiResponse<Product[]>>('/products/best-sellers')
    return response.data.data!
  },

  searchProducts: async (keyword: string): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<Product>>>('/products/search', {
      params: { keyword }
    })
    return response.data.data!
  },

  getProductsByCategory: async (categoryId: number): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<Product>>>(`/products/category/${categoryId}`)
    return response.data.data!
  },

  getCategories: async (): Promise<Category[]> => {
    const response = await apiClient.get<ApiResponse<Category[]>>('/categories')
    return response.data.data!
  },

  getCategoryTree: async (): Promise<Category[]> => {
    const response = await apiClient.get<ApiResponse<Category[]>>('/categories/tree')
    return response.data.data!
  },
}

export default productService
