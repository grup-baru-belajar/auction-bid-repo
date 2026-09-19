import { createSlice, createAsyncThunk, type PayloadAction } from "@reduxjs/toolkit";
import { bidApi } from "../../services/api";
import type { CreateBidRequest, BidResponse } from "../../types";
import { AxiosError } from "axios";

interface BiddingState {
  lastBid: BidResponse | null;
  loading: boolean;
  error: string | null;
}

const initialState: BiddingState = {
  lastBid: null,
  loading: false,
  error: null,
};

// POST /bid
export const placeBid = createAsyncThunk<
  BidResponse,
  CreateBidRequest,
  { rejectValue: string }
>("bidding/placeBid", async (data, { rejectWithValue }) => {
  try {
    const response = await bidApi.createBid(data);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to place bid"
    );
  }
});

const biddingSlice = createSlice({
  name: "bidding",
  initialState,
  reducers: {
    clearBidError: (state) => {
      state.error = null;
    },
    clearLastBid: (state) => {
      state.lastBid = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(placeBid.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(placeBid.fulfilled, (state, action: PayloadAction<BidResponse>) => {
        state.loading = false;
        state.lastBid = action.payload;
      })
      .addCase(placeBid.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Failed to place bid";
      });
  },
});

export const { clearBidError, clearLastBid } = biddingSlice.actions;
export default biddingSlice.reducer;
