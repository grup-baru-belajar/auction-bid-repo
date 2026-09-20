import apiClient from "./apiClient";
import type { ApiResponse, AuctionSummary } from "../../types";

export const reportApi = {
  /** GET /reporting/auction-summary (admin only) */
  getAuctionSummary: () =>
    apiClient.get<ApiResponse<AuctionSummary>>("/v1/reporting/auction-summary"),

  /**
   * GET /reporting/total-bidders (admin only).
   * The endpoint counts distinct bidders within the last `interval` days —
   * there's no "all time" mode, so we pass a large interval to approximate it.
   */
  getTotalBidders: (intervalDays: number) =>
    apiClient.get<ApiResponse<number>>("/v1/reporting/total-bidders", {
      params: { interval: intervalDays },
    }),
};
