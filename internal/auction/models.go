package auction

import "time"

// Auction represents an auction entity.
type Auction struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    StartsAt    time.Time `json:"starts_at"`
    EndsAt      time.Time `json:"ends_at"`
    SellerID    string    `json:"seller_id"`
}

// Bid represents a bid on an auction.
type Bid struct {
    AuctionID string    `json:"auction_id"`
    UserID    string    `json:"user_id"`
    Amount    int64     `json:"amount_cents"`
    PlacedAt  time.Time `json:"placed_at"`
}


