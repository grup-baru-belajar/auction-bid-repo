package models

type TopAuction struct {
	ID            int64   `json:"id"`
	AuctionName   string  `json:"auctionName"`
	StartingPrice float64 `json:"startingPrice"`
	HighestBid    float64 `json:"highestBid"`
	TotalBids     int64   `json:"totalBid"`
	Bidders       int64   `json:"bidders"`
	Status        string  `json:"status"`
}

type AuctionActivity struct {
	Date      string `json:"date"`
	TotalBids int64  `json:"totalBids"`
}

type AuctionStatus struct {
	Status string `json:"status"`
	Total  int64  `json:"total"`
}
