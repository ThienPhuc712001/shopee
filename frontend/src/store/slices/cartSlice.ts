import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Cart, CartItem } from '../../types'

interface CartState {
  cart: Cart | null;
  loading: boolean;
  error: string | null;
}

const initialState: CartState = {
  cart: null,
  loading: false,
  error: null,
}

const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    fetchCartStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    fetchCartSuccess: (state, action: PayloadAction<Cart>) => {
      state.loading = false;
      state.cart = action.payload;
    },
    fetchCartFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    addToCartStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    addToCartSuccess: (state, action: PayloadAction<Cart>) => {
      state.loading = false;
      state.cart = action.payload;
    },
    addToCartFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    updateCartItemQuantity: (state, action: PayloadAction<{ itemId: number; quantity: number }>) => {
      if (state.cart) {
        const item = state.cart.items.find(i => i.id === action.payload.itemId);
        if (item) {
          item.quantity = action.payload.quantity;
          item.subtotal = item.product.price * action.payload.quantity;
          state.cart.subtotal = state.cart.items.reduce(
            (sum, item) => sum + item.subtotal,
            0
          );
        }
      }
    },
    removeFromCart: (state, action: PayloadAction<number>) => {
      if (state.cart) {
        state.cart.items = state.cart.items.filter(i => i.id !== action.payload);
        state.cart.total_items = state.cart.items.reduce((sum, item) => sum + item.quantity, 0);
        state.cart.subtotal = state.cart.items.reduce(
          (sum, item) => sum + item.subtotal,
          0
        );
      }
    },
    clearCart: (state) => {
      state.cart = null;
    },
  },
})

export const {
  fetchCartStart,
  fetchCartSuccess,
  fetchCartFailure,
  addToCartStart,
  addToCartSuccess,
  addToCartFailure,
  updateCartItemQuantity,
  removeFromCart,
  clearCart,
} = cartSlice.actions

export default cartSlice.reducer
