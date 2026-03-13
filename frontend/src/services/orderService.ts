import apiClient from './api'
import { ApiResponse, Order } from '../types'

interface CheckoutData {
  items: Array<{
    product_id: number
    quantity: number
  }>
  shipping_address: {
    name: string
    phone: string
    address: string
    ward?: string
    district?: string
    city: string
  }
  payment_method: string
  coupon_code?: string
}

export const orderService = {
  checkout: async (data: CheckoutData): Promise<Order> => {
    const response = await apiClient.post<ApiResponse<Order>>('/orders', data)
    return response.data.data!
  },

  getOrders: async (): Promise<Order[]> => {
    const response = await apiClient.get<ApiResponse<Order[]>>('/orders')
    return response.data.data!
  },

  getOrderById: async (id: number): Promise<Order> => {
    const response = await apiClient.get<ApiResponse<Order>>(`/orders/${id}`)
    return response.data.data!
  },

  cancelOrder: async (id: number): Promise<Order> => {
    const response = await apiClient.post<ApiResponse<Order>>(`/orders/${id}/cancel`)
    return response.data.data!
  },

  getOrderTracking: async (id: number): Promise<any> => {
    const response = await apiClient.get(`/orders/${id}/tracking`)
    return response.data
  },
}

export default orderService
