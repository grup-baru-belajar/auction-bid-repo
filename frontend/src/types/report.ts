/** One day's worth of activity: auctions created and bids placed. */
export interface AuctionActivity {
  date: string;
  totalAuctions: number;
  totalBids: number;
}

export type AuctionStatusName = "ACTIVE" | "ENDED";

export interface AuctionStatusCount {
  status: AuctionStatusName;
  total: number;
}

/** Sum of completed auctions' lastPrice for one week. */
export interface TransactionWeek {
  weekStart: string;
  total: number;
}

/** Aggregate snapshot returned by GET /reporting/auction-summary */
export interface AuctionSummary {
  totalAuctions: number;
  ongoingAuctions: number;
  completedAuctions: number;
  totalBidsOngoing: number;
  totalBidsCompleted: number;
  totalBidsAll: number;
}
