import apiClient from './api';
import type { Order, OrderStats } from '../types';

export interface CheckoutInput {
  shipping_address_id: number;
  payment_method: string;
  buyer_note?: string;
  shipping_method?: string;
  voucher_code?: string;
}

export interface CreateOrderInput {
  items: {
    product_id: number;
    variant_id?: number;
    quantity: number;
  }[];
  shipping_address: {
    name: string;
    phone: string;
    street: string;
    ward?: string;
    district: string;
    city: string;
    country?: string;
  };
  payment_method: string;
  buyer_note?: string;
}

export const orderService = {
  // Create order from cart
  async checkout(data: CheckoutInput): Promise<Order> {
    const response = await apiClient.post<Order>('/orders/checkout', data);
    return response.data;
  },

  // Create order (direct purchase)
  async createOrder(data: CreateOrderInput): Promise<Order> {
    const response = await apiClient.post<Order>('/orders', data);
    return response.data;
  },

  // Get user orders
  async getUserOrders(page: number = 1, limit: number = 20): Promise<{ orders: Order[]; total: number }> {
    const response = await apiClient.get<{ orders: Order[]; total: number }>(
      `/orders?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Get order by ID
  async getOrderById(id: number): Promise<Order> {
    const response = await apiClient.get<Order>(`/orders/${id}`);
    return response.data;
  },

  // Get order by order number
  async getOrderByOrderNumber(orderNumber: string): Promise<Order> {
    const response = await apiClient.get<Order>(`/orders/${orderNumber}`);
    return response.data;
  },

  // Cancel order
  async cancelOrder(id: number, reason: string): Promise<Order> {
    const response = await apiClient.post<Order>(`/orders/${id}/cancel`, { reason });
    return response.data;
  },

  // Confirm order received
  async confirmOrder(id: number): Promise<Order> {
    const response = await apiClient.post<Order>(`/orders/${id}/confirm`);
    return response.data;
  },

  // Request refund
  async requestRefund(id: number, reason: string, amount?: number): Promise<Order> {
    const response = await apiClient.post<Order>(`/orders/${id}/refund`, { reason, amount });
    return response.data;
  },

  // Get order tracking
  async getOrderTracking(id: number): Promise<{ status: string; events: { status: string; description: string; timestamp: string }[] }> {
    const response = await apiClient.get(`/orders/${id}/tracking`);
    return response.data as { status: string; events: { status: string; description: string; timestamp: string }[] };
  },

  // Get order statistics
  async getOrderStats(): Promise<OrderStats> {
    const response = await apiClient.get<OrderStats>('/orders/stats');
    return response.data;
  },

  // Get orders by status
  async getOrdersByStatus(status: string, page: number = 1, limit: number = 20): Promise<{ orders: Order[]; total: number }> {
    const response = await apiClient.get<{ orders: Order[]; total: number }>(
      `/orders?status=${status}&page=${page}&limit=${limit}`
    );
    return response.data;
  },
};

export default orderService;
