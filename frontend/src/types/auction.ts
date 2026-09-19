export interface BidWinner {
  id: number;
  name: string;
}

export interface TopBid {
  id: number;
  userId: number;
  userName: string;
  bidPrice: number;
  createdAt: string;
}

/** Auction item as returned by GET /auctions (list) */
export interface Auction {
  id: number;
  auctionName: string;
  description: string;
  imageLink: string;
  startingPrice: number;
  lastPrice: number;
  createdAt: string;
  endTime: string;
  isCompleted: boolean;
  bidWinner: BidWinner | null;
}

/** Auction detail as returned by GET /auctions/{id} — includes topBids */
export interface AuctionDetail extends Auction {
  topBids: TopBid[];
  totalBids: number;
  totalBidders: number;

}

/** Request body for POST /auctions */
export interface CreateAuctionRequest {
  auctionName: string;
  description: string;
  imageLink: string;
  startingPrice: number;
  endTime: string;
}

/** Query params for GET /auctions */
export interface GetAuctionsParams {
  page?: number;
  limit?: number;
  isCompleted?: boolean;
}
