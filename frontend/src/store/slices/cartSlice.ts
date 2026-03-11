import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Cart, CartItem } from '../../types';
import cartService from '../../services/cartService';

interface CartState {
  cart: Cart | null;
  items: CartItem[];
  isLoading: boolean;
  error: string | null;
  itemCount: number;
}

const initialState: CartState = {
  cart: null,
  items: [],
  isLoading: false,
  error: null,
  itemCount: 0,
};

export const fetchCart = createAsyncThunk('cart/fetchCart', async (_, { rejectWithValue }) => {
  try {
    const cart = await cartService.getCart();
    return cart;
  } catch (error: unknown) {
    const err = error as { message: string };
    return rejectWithValue(err.message);
  }
});

export const addToCart = createAsyncThunk(
  'cart/addToCart',
  async ({ productId, quantity, variantId }: { productId: number; quantity: number; variantId?: number }, { rejectWithValue }) => {
    try {
      const cart = await cartService.addItem(productId, quantity, variantId);
      return cart;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const updateCartItem = createAsyncThunk(
  'cart/updateCartItem',
  async ({ itemId, quantity }: { itemId: number; quantity: number }, { rejectWithValue }) => {
    try {
      const cart = await cartService.updateItem(itemId, quantity);
      return cart;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const removeFromCart = createAsyncThunk(
  'cart/removeFromCart',
  async (itemId: number, { rejectWithValue }) => {
    try {
      const cart = await cartService.removeItem(itemId);
      return cart;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const clearCart = createAsyncThunk('cart/clearCart', async (_, { rejectWithValue }) => {
  try {
    const cart = await cartService.clearCart();
    return cart;
  } catch (error: unknown) {
    const err = error as { message: string };
    return rejectWithValue(err.message);
  }
});

const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    clearCartError: (state) => {
      state.error = null;
    },
    setItemCount: (state, action) => {
      state.itemCount = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch cart
      .addCase(fetchCart.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchCart.fulfilled, (state, action) => {
        state.isLoading = false;
        state.cart = action.payload;
        state.items = action.payload.items;
        state.itemCount = action.payload.total_items;
      })
      .addCase(fetchCart.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Add to cart
      .addCase(addToCart.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(addToCart.fulfilled, (state, action) => {
        state.isLoading = false;
        state.cart = action.payload;
        state.items = action.payload.items;
        state.itemCount = action.payload.total_items;
      })
      .addCase(addToCart.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Update cart item
      .addCase(updateCartItem.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(updateCartItem.fulfilled, (state, action) => {
        state.isLoading = false;
        state.cart = action.payload;
        state.items = action.payload.items;
        state.itemCount = action.payload.total_items;
      })
      .addCase(updateCartItem.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Remove from cart
      .addCase(removeFromCart.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(removeFromCart.fulfilled, (state, action) => {
        state.isLoading = false;
        state.cart = action.payload;
        state.items = action.payload.items;
        state.itemCount = action.payload.total_items;
      })
      .addCase(removeFromCart.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Clear cart
      .addCase(clearCart.fulfilled, (state, action) => {
        state.cart = action.payload;
        state.items = action.payload.items;
        state.itemCount = action.payload.total_items;
      });
  },
});

export const { clearCartError, setItemCount } = cartSlice.actions;
export default cartSlice.reducer;
