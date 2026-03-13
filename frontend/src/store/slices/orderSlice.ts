import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Order } from '../../types'

interface OrderState {
  orders: Order[];
  selectedOrder: Order | null;
  loading: boolean;
  error: string | null;
}

const initialState: OrderState = {
  orders: [],
  selectedOrder: null,
  loading: false,
  error: null,
}

const orderSlice = createSlice({
  name: 'orders',
  initialState,
  reducers: {
    fetchOrdersStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    fetchOrdersSuccess: (state, action: PayloadAction<Order[]>) => {
      state.loading = false;
      state.orders = action.payload;
    },
    fetchOrdersFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    fetchOrderByIdStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    fetchOrderByIdSuccess: (state, action: PayloadAction<Order>) => {
      state.loading = false;
      state.selectedOrder = action.payload;
    },
    fetchOrderByIdFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    createOrderStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    createOrderSuccess: (state, action: PayloadAction<Order>) => {
      state.loading = false;
      state.orders.unshift(action.payload);
    },
    createOrderFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    cancelOrderSuccess: (state, action: PayloadAction<number>) => {
      const order = state.orders.find(o => o.id === action.payload);
      if (order) {
        order.status = 'cancelled';
      }
    },
  },
})

export const {
  fetchOrdersStart,
  fetchOrdersSuccess,
  fetchOrdersFailure,
  fetchOrderByIdStart,
  fetchOrderByIdSuccess,
  fetchOrderByIdFailure,
  createOrderStart,
  createOrderSuccess,
  createOrderFailure,
  cancelOrderSuccess,
} = orderSlice.actions

export default orderSlice.reducer
