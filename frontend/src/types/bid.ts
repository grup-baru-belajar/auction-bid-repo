/** Request body for POST /bid */
export interface CreateBidRequest {
  auctionId: number;
  bidPrice: number;
}

/** Response data from POST /bid */
export interface BidResponse {
  id: number;
  auctionId: number;
  userId: number;
  bidPrice: number;
  createdAt: string;
}
