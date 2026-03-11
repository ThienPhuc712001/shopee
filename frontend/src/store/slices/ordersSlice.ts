import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Order, OrderStats } from '../../types';
import orderService from '../../services/orderService';

interface OrdersState {
  orders: Order[];
  currentOrder: Order | null;
  stats: OrderStats | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
  isLoading: boolean;
  isCheckingOut: boolean;
  error: string | null;
}

const initialState: OrdersState = {
  orders: [],
  currentOrder: null,
  stats: null,
  pagination: {
    page: 1,
    limit: 20,
    total: 0,
  },
  isLoading: false,
  isCheckingOut: false,
  error: null,
};

export const fetchUserOrders = createAsyncThunk(
  'orders/fetchUserOrders',
  async ({ page, limit }: { page: number; limit: number }, { rejectWithValue }) => {
    try {
      const data = await orderService.getUserOrders(page, limit);
      return data;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const fetchOrderById = createAsyncThunk(
  'orders/fetchOrderById',
  async (id: number, { rejectWithValue }) => {
    try {
      const order = await orderService.getOrderById(id);
      return order;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const checkout = createAsyncThunk(
  'orders/checkout',
  async (data: { shipping_address_id: number; payment_method: string; buyer_note?: string; voucher_code?: string }, { rejectWithValue }) => {
    try {
      const order = await orderService.checkout(data);
      return order;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const cancelOrder = createAsyncThunk(
  'orders/cancelOrder',
  async ({ id, reason }: { id: number; reason: string }, { rejectWithValue }) => {
    try {
      const order = await orderService.cancelOrder(id, reason);
      return order;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const fetchOrderStats = createAsyncThunk(
  'orders/fetchOrderStats',
  async (_, { rejectWithValue }) => {
    try {
      const stats = await orderService.getOrderStats();
      return stats;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

const ordersSlice = createSlice({
  name: 'orders',
  initialState,
  reducers: {
    clearCurrentOrder: (state) => {
      state.currentOrder = null;
    },
    clearOrdersError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch user orders
      .addCase(fetchUserOrders.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUserOrders.fulfilled, (state, action) => {
        state.isLoading = false;
        state.orders = action.payload.orders;
        state.pagination.total = action.payload.total;
      })
      .addCase(fetchUserOrders.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch order by ID
      .addCase(fetchOrderById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchOrderById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentOrder = action.payload;
      })
      .addCase(fetchOrderById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Checkout
      .addCase(checkout.pending, (state) => {
        state.isCheckingOut = true;
        state.error = null;
      })
      .addCase(checkout.fulfilled, (state, action) => {
        state.isCheckingOut = false;
        state.currentOrder = action.payload;
      })
      .addCase(checkout.rejected, (state, action) => {
        state.isCheckingOut = false;
        state.error = action.payload as string;
      })
      // Cancel order
      .addCase(cancelOrder.fulfilled, (state, action) => {
        const index = state.orders.findIndex((o) => o.id === action.payload.id);
        if (index !== -1) {
          state.orders[index] = action.payload;
        }
        if (state.currentOrder?.id === action.payload.id) {
          state.currentOrder = action.payload;
        }
      })
      // Fetch order stats
      .addCase(fetchOrderStats.fulfilled, (state, action) => {
        state.stats = action.payload;
      });
  },
});

export const { clearCurrentOrder, clearOrdersError } = ordersSlice.actions;
export default ordersSlice.reducer;
