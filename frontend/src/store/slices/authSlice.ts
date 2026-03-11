import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import type { User } from '../../types';
import authService from '../../services/authService';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

const storedUser = authService.getStoredUser();
const storedToken = localStorage.getItem('access_token');

const initialState: AuthState = {
  user: storedUser,
  token: storedToken,
  isAuthenticated: !!storedToken,
  isLoading: false,
  error: null,
};

export const login = createAsyncThunk(
  'auth/login',
  async (credentials: { email: string; password: string }, { rejectWithValue }) => {
    try {
      const data = await authService.login(credentials);
      return data;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const register = createAsyncThunk(
  'auth/register',
  async (data: { email: string; password: string; first_name: string; last_name: string; phone: string }, { rejectWithValue }) => {
    try {
      const response = await authService.register(data);
      return response;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const logout = createAsyncThunk(
  'auth/logout',
  async (_, { rejectWithValue }) => {
    try {
      await authService.logout();
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

export const getCurrentUser = createAsyncThunk(
  'auth/getCurrentUser',
  async (_, { rejectWithValue }) => {
    try {
      const user = await authService.getCurrentUser();
      return user;
    } catch (error: unknown) {
      const err = error as { message: string };
      return rejectWithValue(err.message);
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    updateUser: (state, action: PayloadAction<Partial<User>>) => {
      if (state.user) {
        state.user = { ...state.user, ...action.payload };
        localStorage.setItem('user', JSON.stringify(state.user));
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Login
      .addCase(login.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.token = action.payload.token;
      })
      .addCase(login.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Register
      .addCase(register.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.token = action.payload.token;
      })
      .addCase(register.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Logout
      .addCase(logout.fulfilled, (state) => {
        state.isAuthenticated = false;
        state.user = null;
        state.token = null;
      })
      // Get Current User
      .addCase(getCurrentUser.fulfilled, (state, action) => {
        state.user = action.payload;
      })
      .addCase(getCurrentUser.rejected, (state) => {
        state.isAuthenticated = false;
        state.user = null;
        state.token = null;
      });
  },
});

export const { clearError, updateUser } = authSlice.actions;
export default authSlice.reducer;
