import apiClient from './api'
import { ApiResponse, Cart, CartItem } from '../types'

export const cartService = {
  getCart: async (): Promise<Cart> => {
    const response = await apiClient.get<ApiResponse<Cart>>('/cart')
    return response.data.data!
  },

  addToCart: async (product_id: number, quantity: number): Promise<Cart> => {
    const response = await apiClient.post<ApiResponse<Cart>>('/cart/add', { product_id, quantity })
    return response.data.data!
  },

  updateCartItem: async (itemId: number, quantity: number): Promise<Cart> => {
    const response = await apiClient.put<ApiResponse<Cart>>(`/cart/items/${itemId}`, { quantity })
    return response.data.data!
  },

  removeFromCart: async (itemId: number): Promise<void> => {
    await apiClient.delete(`/cart/items/${itemId}`)
  },

  clearCart: async (): Promise<void> => {
    await apiClient.delete('/cart/clear')
  },
}

export default cartService
