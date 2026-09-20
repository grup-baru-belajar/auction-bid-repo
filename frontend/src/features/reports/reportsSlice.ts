import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { reportApi } from "../../services/api";
import type {
  AuctionSummary,
  AuctionActivity,
  AuctionStatusCount,
  TransactionWeek,
} from "../../types";
import { AxiosError } from "axios";

/** No "all time" mode on the backend, so we ask for a wide window instead. */
// const ALL_TIME_INTERVAL_DAYS = 36500;

interface ReportsState {
  summary: AuctionSummary | null;
  loading: boolean;
  error: string | null;

  totalBidders: number | null;
  biddersError: string | null;

  activity: AuctionActivity[] | null;
  activityIntervalDays: number;
  activityLoading: boolean;
  activityError: string | null;

  statusBreakdown: AuctionStatusCount[] | null;
  statusError: string | null;

  transactionOverview: TransactionWeek[] | null;
  transactionError: string | null;
}

const initialState: ReportsState = {
  summary: null,
  loading: false,
  error: null,

  totalBidders: null,
  biddersError: null,

  activity: null,
  activityIntervalDays: 7,
  activityLoading: false,
  activityError: null,

  statusBreakdown: null,
  statusError: null,

  transactionOverview: null,
  transactionError: null,
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
    const response = await reportApi.getTotalBidders("all");
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch total bidders",
    );
  }
});

// GET /reporting/auction-activity
export const fetchAuctionActivity = createAsyncThunk<
  AuctionActivity[],
  number,
  { rejectValue: string }
>("reports/fetchAuctionActivity", async (intervalDays, { rejectWithValue }) => {
  try {
    const response = await reportApi.getAuctionActivity(intervalDays);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch auction activity",
    );
  }
});

// GET /reporting/auction-status
export const fetchAuctionStatus = createAsyncThunk<
  AuctionStatusCount[],
  void,
  { rejectValue: string }
>("reports/fetchAuctionStatus", async (_, { rejectWithValue }) => {
  try {
    const response = await reportApi.getAuctionStatus();
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch auction status",
    );
  }
});

// GET /reporting/transaction-overview
export const fetchTransactionOverview = createAsyncThunk<
  TransactionWeek[],
  number,
  { rejectValue: string }
>("reports/fetchTransactionOverview", async (weeks, { rejectWithValue }) => {
  try {
    const response = await reportApi.getTransactionOverview(weeks);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch transaction overview",
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
    setActivityIntervalDays: (state, action: { payload: number }) => {
      state.activityIntervalDays = action.payload;
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
      })
      .addCase(fetchAuctionActivity.pending, (state) => {
        state.activityLoading = true;
        state.activityError = null;
      })
      .addCase(fetchAuctionActivity.fulfilled, (state, action) => {
        state.activityLoading = false;
        state.activity = action.payload;
      })
      .addCase(fetchAuctionActivity.rejected, (state, action) => {
        state.activityLoading = false;
        state.activityError =
          action.payload || "Failed to fetch auction activity";
      })
      .addCase(fetchAuctionStatus.fulfilled, (state, action) => {
        state.statusBreakdown = action.payload;
      })
      .addCase(fetchAuctionStatus.rejected, (state, action) => {
        state.statusError = action.payload || "Failed to fetch auction status";
      })
      .addCase(fetchTransactionOverview.fulfilled, (state, action) => {
        state.transactionOverview = action.payload;
      })
      .addCase(fetchTransactionOverview.rejected, (state, action) => {
        state.transactionError =
          action.payload || "Failed to fetch transaction overview";
      });
  },
});

export const { clearError, setActivityIntervalDays } = reportsSlice.actions;
export default reportsSlice.reducer;
