import apiClient from "./apiClient";
import type {
  ApiResponse,
  AuctionSummary,
  AuctionActivity,
  AuctionStatusCount,
  TransactionWeek,
  TopAuction,
} from "../../types";

export const reportApi = {
  /** GET /reporting/auction-summary (admin only) */
  getAuctionSummary: () =>
    apiClient.get<ApiResponse<AuctionSummary>>("/v1/reporting/auction-summary"),

  /** GET /reporting/total-bidders (admin only). interval="all" counts all-time. */
  getTotalBidders: (interval: string | number = "all") =>
    apiClient.get<ApiResponse<number>>("/v1/reporting/total-bidders", {
      params: { interval },
    }),

  /** GET /reporting/auction-activity (admin only). Auctions created + bids placed per day. */
  getAuctionActivity: (intervalDays: number) =>
    apiClient.get<ApiResponse<AuctionActivity[]>>(
      "/v1/reporting/auction-activity",
      {
        params: { interval: intervalDays },
      },
    ),

  /** GET /reporting/auction-status (admin only). ACTIVE/ENDED counts. */
  getAuctionStatus: () =>
    apiClient.get<ApiResponse<AuctionStatusCount[]>>(
      "/v1/reporting/auction-status",
    ),

  /** GET /reporting/transaction-overview (admin only). Completed-auction value per week. */
  getTransactionOverview: (weeks: number) =>
    apiClient.get<ApiResponse<TransactionWeek[]>>(
      "/v1/reporting/transaction-overview",
      {
        params: { weeks },
      },
    ),

  /** GET /reporting/top-auction (admin only). Top 5 auctions by total bids. */
  getTopAuctions: () =>
    apiClient.get<ApiResponse<TopAuction[]>>("/v1/reporting/top-auction"),
};
