import {
  createSlice,
  createAsyncThunk,
  type PayloadAction,
} from "@reduxjs/toolkit";
import { auctionApi } from "../../services/api";
import type {
  Auction,
  AuctionDetail,
  CreateAuctionRequest,
  GetAuctionsParams,
  Pagination,
} from "../../types";
import { AxiosError } from "axios";

interface AuctionsState {
  list: Auction[];
  selectedAuction: AuctionDetail | null;
  pagination: Pagination | null;
  loading: boolean;
  detailLoading: boolean;
  error: string | null;
}

const initialState: AuctionsState = {
  list: [],
  selectedAuction: null,
  pagination: null,
  loading: false,
  detailLoading: false,
  error: null,
};

// GET /auctions
export const fetchAuctions = createAsyncThunk<
  { data: Auction[]; pagination: Pagination },
  GetAuctionsParams | undefined,
  { rejectValue: string }
>("auctions/fetchAuctions", async (params, { rejectWithValue }) => {
  try {
    const response = await auctionApi.getAuctions(params);
    return {
      data: response.data.data,
      pagination: response.data.pagination,
    };
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch auctions",
    );
  }
});

// GET /auctions/{id}
export const fetchAuctionDetail = createAsyncThunk<
  AuctionDetail,
  number,
  { rejectValue: string }
>("auctions/fetchAuctionDetail", async (id, { rejectWithValue }) => {
  try {
    const response = await auctionApi.getAuctionDetail(id);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to fetch auction detail",
    );
  }
});

// POST /auctions
export const createAuction = createAsyncThunk<
  Auction,
  CreateAuctionRequest,
  { rejectValue: string }
>("auctions/createAuction", async (data, { rejectWithValue }) => {
  try {
    const response = await auctionApi.createAuction(data);
    return response.data.data;
  } catch (err) {
    const error = err as AxiosError<{ message: string }>;
    return rejectWithValue(
      error.response?.data?.message || "Failed to create auction",
    );
  }
});

const auctionsSlice = createSlice({
  name: "auctions",
  initialState,
  reducers: {
    clearSelectedAuction: (state) => {
      state.selectedAuction = null;
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // fetchAuctions
      .addCase(fetchAuctions.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAuctions.fulfilled, (state, action) => {
        state.loading = false;
        state.list = action.payload.data;
        state.pagination = action.payload.pagination;
      })
      .addCase(fetchAuctions.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Failed to fetch auctions";
      })
      // fetchAuctionDetail
      .addCase(fetchAuctionDetail.pending, (state) => {
        state.detailLoading = true;
        state.error = null;
      })
      .addCase(
        fetchAuctionDetail.fulfilled,
        (state, action: PayloadAction<AuctionDetail>) => {
          state.detailLoading = false;
          state.selectedAuction = action.payload;
        },
      )
      .addCase(fetchAuctionDetail.rejected, (state, action) => {
        state.detailLoading = false;
        state.error = action.payload || "Failed to fetch auction detail";
      })
      // createAuction
      .addCase(createAuction.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        createAuction.fulfilled,
        (state, action: PayloadAction<Auction>) => {
          state.loading = false;
          state.list.unshift(action.payload);
        },
      )
      .addCase(createAuction.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Failed to create auction";
      });
  },
});

export const { clearSelectedAuction, clearError } = auctionsSlice.actions;
export default auctionsSlice.reducer;
