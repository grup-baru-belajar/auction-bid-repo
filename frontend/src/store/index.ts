import { configureStore } from "@reduxjs/toolkit";
import authReducer from "../features/auth/authSlice";
import auctionsReducer from "../features/auctions/auctionsSlice";
import biddingReducer from "../features/bidding/biddingSlice";
import reportsReducer from "../features/reports/reportsSlice";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    auctions: auctionsReducer,
    bidding: biddingReducer,
    reports: reportsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
