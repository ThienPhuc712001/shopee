import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Product, Category } from '../../types'

interface ProductState {
  products: Product[];
  categories: Category[];
  selectedProduct: Product | null;
  loading: boolean;
  error: string | null;
  filters: {
    category_id?: number;
    min_price?: number;
    max_price?: number;
    keyword?: string;
    sort_by: string;
    sort_order: 'asc' | 'desc';
  };
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

const initialState: ProductState = {
  products: [],
  categories: [],
  selectedProduct: null,
  loading: false,
  error: null,
  filters: {
    sort_by: 'created_at',
    sort_order: 'desc',
  },
  pagination: {
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  },
}

const productSlice = createSlice({
  name: 'products',
  initialState,
  reducers: {
    fetchProductsStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    fetchProductsSuccess: (state, action: PayloadAction<{ products: Product[]; pagination: any }>) => {
      state.loading = false;
      state.products = action.payload.products;
      state.pagination = action.payload.pagination;
    },
    fetchProductsFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    fetchProductByIdStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    fetchProductByIdSuccess: (state, action: PayloadAction<Product>) => {
      state.loading = false;
      state.selectedProduct = action.payload;
    },
    fetchProductByIdFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    fetchCategoriesSuccess: (state, action: PayloadAction<Category[]>) => {
      state.categories = action.payload;
    },
    setFilters: (state, action: PayloadAction<Partial<typeof initialState.filters>>) => {
      state.filters = { ...state.filters, ...action.payload };
      state.pagination.page = 1;
    },
    resetFilters: (state) => {
      state.filters = initialState.filters;
      state.pagination.page = 1;
    },
  },
})

export const {
  fetchProductsStart,
  fetchProductsSuccess,
  fetchProductsFailure,
  fetchProductByIdStart,
  fetchProductByIdSuccess,
  fetchProductByIdFailure,
  fetchCategoriesSuccess,
  setFilters,
  resetFilters,
} = productSlice.actions

export default productSlice.reducer
