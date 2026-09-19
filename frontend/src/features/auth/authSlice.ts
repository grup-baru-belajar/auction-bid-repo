import { createSlice, createAsyncThunk, type PayloadAction } from "@reduxjs/toolkit";
import { authApi } from "../../services/api";
import type { LoginRequest, LoginResponse, Role } from "../../types";
import { AxiosError } from "axios";

interface AuthState {
  user: {
    id: number;
    name: string;
    username: string;
    role: Role;
  } | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

const storedToken = localStorage.getItem("token");
const storedUser = localStorage.getItem("user");

const initialState: AuthState = {
  user: storedUser ? JSON.parse(storedUser) : null,
  token: storedToken,
  isAuthenticated: !!storedToken,
  loading: false,
  error: null,
};

// Async thunk — login
export const login = createAsyncThunk<
  LoginResponse,
  LoginRequest,
  { rejectValue: string }
>("auth/login", async (credentials, { rejectWithValue }) => {
  try {
    const response = await authApi.login(credentials);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Login failed"
    );
  }
});

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.error = null;
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action: PayloadAction<LoginResponse>) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.token = action.payload.token;
        state.user = {
          id: action.payload.id,
          name: action.payload.name,
          username: action.payload.username,
          role: action.payload.role,
        };
        localStorage.setItem("token", action.payload.token);
        localStorage.setItem("user", JSON.stringify(state.user));
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Login failed";
      });
  },
});

export const { logout, clearError } = authSlice.actions;
export default authSlice.reducer;
