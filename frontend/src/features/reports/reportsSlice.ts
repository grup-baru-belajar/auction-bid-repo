import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { reportApi } from "../../services/api";
import type { AuctionSummary } from "../../types";
import { AxiosError } from "axios";

/** No "all time" mode on the backend, so we ask for a wide window instead. */
const ALL_TIME_INTERVAL_DAYS = 36500;

interface ReportsState {
  summary: AuctionSummary | null;
  loading: boolean;
  error: string | null;
  totalBidders: number | null;
  biddersError: string | null;
}

const initialState: ReportsState = {
  summary: null,
  loading: false,
  error: null,
  totalBidders: null,
  biddersError: null,
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

// GET /reporting/total-bidders
export const fetchTotalBidders = createAsyncThunk<
  number,
  void,
  { rejectValue: string }
>("reports/fetchTotalBidders", async (_, { rejectWithValue }) => {
  try {
    const response = await reportApi.getTotalBidders(ALL_TIME_INTERVAL_DAYS);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch total bidders",
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
      })
      .addCase(fetchTotalBidders.fulfilled, (state, action) => {
        state.totalBidders = action.payload;
      })
      .addCase(fetchTotalBidders.rejected, (state, action) => {
        state.biddersError = action.payload || "Failed to fetch total bidders";
      });
  },
});

export const { clearError } = reportsSlice.actions;
export default reportsSlice.reducer;
