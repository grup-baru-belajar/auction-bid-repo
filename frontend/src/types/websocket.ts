import type { TopBid } from "./auction";

export interface TopBidsWsMessage {
  auctionId: number;
  totalBids: number;
  totalBidders: number;
  topBids: TopBid[];
}
