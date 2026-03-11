import apiClient from './api';
import type { Cart } from '../types';

export const cartService = {
  // Get cart
  async getCart(): Promise<Cart> {
    const response = await apiClient.get<Cart>('/cart');
    return response.data;
  },

  // Add item to cart
  async addItem(productId: number, quantity: number, variantId?: number): Promise<Cart> {
    const response = await apiClient.post<Cart>('/cart/items', {
      product_id: productId,
      quantity,
      variant_id: variantId,
    });
    return response.data;
  },

  // Update cart item quantity
  async updateItem(itemId: number, quantity: number): Promise<Cart> {
    const response = await apiClient.put<Cart>(`/cart/items/${itemId}`, { quantity });
    return response.data;
  },

  // Remove item from cart
  async removeItem(itemId: number): Promise<Cart> {
    const response = await apiClient.delete<Cart>(`/cart/items/${itemId}`);
    return response.data;
  },

  // Clear cart
  async clearCart(): Promise<Cart> {
    const response = await apiClient.delete<Cart>('/cart/clear');
    return response.data;
  },

  // Apply discount code
  async applyDiscount(code: string): Promise<Cart> {
    const response = await apiClient.post<Cart>('/cart/discount', { code });
    return response.data;
  },

  // Remove discount
  async removeDiscount(): Promise<Cart> {
    const response = await apiClient.delete<Cart>('/cart/discount');
    return response.data;
  },

  // Get cart item count
  async getItemCount(): Promise<number> {
    const cart = await this.getCart();
    return cart.total_items;
  },
};

export default cartService;
