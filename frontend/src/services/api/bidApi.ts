import apiClient from "./apiClient";
import type { ApiResponse, CreateBidRequest, BidResponse } from "../../types";

export const bidApi = {
  /** POST /bid — place a bid on an auction */
  createBid: (data: CreateBidRequest) =>
    apiClient.post<ApiResponse<BidResponse>>("/bid", data),
};
