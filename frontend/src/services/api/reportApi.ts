import apiClient from "./apiClient";
import type { ApiResponse, AuctionSummary } from "../../types";

export const reportApi = {
  /** GET /reporting/auction-summary (admin only) */
  getAuctionSummary: () =>
    apiClient.get<ApiResponse<AuctionSummary>>("/v1/reporting/auction-summary"),
};
