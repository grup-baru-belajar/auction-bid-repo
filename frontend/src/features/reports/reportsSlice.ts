import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { reportApi } from "../../services/api";
import type { AuctionSummary } from "../../types";
import { AxiosError } from "axios";

interface ReportsState {
  summary: AuctionSummary | null;
  loading: boolean;
  error: string | null;
}

const initialState: ReportsState = {
  summary: null,
  loading: false,
  error: null,
};

// GET /reporting/auction-summary
export const fetchAuctionSummary = createAsyncThunk<
  AuctionSummary,
  void,
  { rejectValue: string }
>("reports/fetchAuctionSummary", async (_, { rejectWithValue }) => {
  try {
    const response = await reportApi.getAuctionSummary();
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch auction summary",
    );
  }
});

const reportsSlice = createSlice({
  name: "reports",
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchAuctionSummary.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAuctionSummary.fulfilled, (state, action) => {
        state.loading = false;
        state.summary = action.payload;
      })
      .addCase(fetchAuctionSummary.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Failed to fetch auction summary";
      });
  },
});

export const { clearError } = reportsSlice.actions;
export default reportsSlice.reducer;
