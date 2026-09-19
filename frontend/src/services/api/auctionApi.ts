import apiClient from "./apiClient";
import type {
  ApiResponse,
  PaginatedResponse,
  Auction,
  AuctionDetail,
  CreateAuctionRequest,
  GetAuctionsParams,
} from "../../types";

export const auctionApi = {
  /** GET /auctions — list with pagination & optional filter */
  getAuctions: (params?: GetAuctionsParams) =>
    apiClient.get<PaginatedResponse<Auction>>("v1/auctions", { params }),

  /** GET /auctions/{id} — detail with topBids */
  getAuctionDetail: (id: number) =>
    apiClient.get<ApiResponse<AuctionDetail>>(`/v1/auctions/${id}`),

  /** POST /auctions — admin only */
  createAuction: (data: CreateAuctionRequest) =>
    apiClient.post<ApiResponse<Auction>>("/v1/auctions", data),
};
