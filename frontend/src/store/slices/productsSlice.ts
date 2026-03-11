import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Product, ProductFilters, Category, PaginatedResponse } from '../../types';
import productService from '../../services/productService';

interface ProductsState {
  products: Product[];
  featuredProducts: Product[];
  flashSaleProducts: Product[];
  categories: Category[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
  filters: ProductFilters;
  isLoading: boolean;
  error: string | null;
}

const initialState: ProductsState = {
  products: [],
  featuredProducts: [],
  flashSaleProducts: [],
  categories: [],
  pagination: {
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  },
  filters: {},
  isLoading: false,
  error: null,
};

export const fetchProducts = createAsyncThunk(
  'products/fetchProducts',
  async (filters: ProductFilters, { rejectWithValue }) => {
    try {
      const data = await productService.getProducts(filters);
      return data;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const fetchFeaturedProducts = createAsyncThunk(
  'products/fetchFeaturedProducts',
  async (limit: number = 10, { rejectWithValue }) => {
    try {
      const products = await productService.getFeaturedProducts(limit);
      return products;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const fetchFlashSaleProducts = createAsyncThunk(
  'products/fetchFlashSaleProducts',
  async (_, { rejectWithValue }) => {
    try {
      const products = await productService.getFlashSaleProducts();
      return products;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const fetchCategories = createAsyncThunk(
  'products/fetchCategories',
  async (_, { rejectWithValue }) => {
    try {
      const categories = await productService.getCategories();
      return categories;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const fetchProductById = createAsyncThunk(
  'products/fetchProductById',
  async (id: number, { rejectWithValue }) => {
    try {
      const product = await productService.getProductById(id);
      return product;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

const productsSlice = createSlice({
  name: 'products',
  initialState,
  reducers: {
    setFilters: (state, action) => {
      state.filters = { ...state.filters, ...action.payload };
      state.pagination.page = 1;
    },
    setPage: (state, action) => {
      state.pagination.page = action.payload;
    },
    clearProductsError: (state) => {
      state.error = null;
    },
    resetProducts: (state) => {
      state.products = [];
      state.pagination = initialState.pagination;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch products
      .addCase(fetchProducts.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchProducts.fulfilled, (state, action) => {
        state.isLoading = false;
        state.products = action.payload.items;
        state.pagination = {
          page: action.payload.page,
          limit: action.payload.limit,
          total: action.payload.total,
          total_pages: action.payload.total_pages,
        };
      })
      .addCase(fetchProducts.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch featured products
      .addCase(fetchFeaturedProducts.fulfilled, (state, action) => {
        state.featuredProducts = action.payload;
      })
      // Fetch flash sale products
      .addCase(fetchFlashSaleProducts.fulfilled, (state, action) => {
        state.flashSaleProducts = action.payload;
      })
      // Fetch categories
      .addCase(fetchCategories.fulfilled, (state, action) => {
        state.categories = action.payload;
      })
      .addCase(fetchCategories.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(fetchCategories.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { setFilters, setPage, clearProductsError, resetProducts } = productsSlice.actions;
export default productsSlice.reducer;
