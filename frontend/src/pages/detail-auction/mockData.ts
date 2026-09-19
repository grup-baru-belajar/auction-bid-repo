import type { AuctionDetail } from "../../types";

export const mockAuctionDetail: AuctionDetail = {
  id: 1,
  auctionName: "Iphone 17 Shahan",
  description: "This is a mock auction",
  imageLink: "https://cdn.example.com/images/iphone17.jpg",
  startingPrice: 999999999,
  lastPrice: 27300000,
  totalBids: 1284,
  totalBidders: 342,
  createdAt: "2024-12-07T09:00:00Z",
  endTime: "2024-12-10T09:00:00Z",
  isCompleted: true,
  bidWinner: {
    id: 1,
    name: "Surya A.",
  },
  topBids: [
    {
      id: 1,
      userId: 101,
      userName: "Surya A.",
      bidPrice: 24750000,
      createdAt: "2024-12-08T10:15:00Z",
    },
    {
      id: 2,
      userId: 102,
      userName: "Surya B.",
      bidPrice: 27300000,
      createdAt: "2024-12-09T14:30:00Z",
    },
    {
      id: 3,
      userId: 103,
      userName: "Surya C.",
      bidPrice: 9850000,
      createdAt: "2024-12-07T09:00:00Z",
    },
  ],
};

/** Simulate an async API call with a small delay */
export const fetchMockAuctionDetail = (
  _id: number
): Promise<AuctionDetail> => {
  return new Promise((resolve) => {
    setTimeout(() => resolve(mockAuctionDetail), 600);
  });
};
