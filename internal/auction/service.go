package auction

import (
    "context"
    "errors"
    "fmt"
    "time"

    ihz "golang-sba-hazelcast/internal/platform/hazelcast"
)

// Service provides core auction operations using Hazelcast as the space.
type Service struct {
    space ihz.Space
}

// NewService creates a new auction service.
func NewService(space ihz.Space) *Service {
    return &Service{space: space}
}

// PlaceBid places a bid with a simple optimistic strategy.
func (s *Service) PlaceBid(ctx context.Context, bid Bid) error {
    if bid.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    m, err := s.space.GetMap(ctx, fmt.Sprintf("auction:%s:bids", bid.AuctionID))
    if err != nil {
        return err
    }
    cctx, cancel := ihz.WithTimeout(ctx)
    defer cancel()
    key := fmt.Sprintf("%s:%d", bid.UserID, bid.PlacedAt.UnixNano())
    _, err = m.Put(cctx, key, bid)
    return err
}

// CreateAuction stores an auction in the distributed map.
func (s *Service) CreateAuction(ctx context.Context, a Auction) error {
    if a.ID == "" || a.Title == "" || a.SellerID == "" {
        return errors.New("missing required fields")
    }
    if !a.StartsAt.Before(a.EndsAt) {
        return errors.New("invalid time window")
    }
    m, err := s.space.GetMap(ctx, "auctions")
    if err != nil {
        return err
    }
    cctx, cancel := ihz.WithTimeout(ctx)
    defer cancel()
    _, err = m.Put(cctx, a.ID, a)
    return err
}

// ActiveAuctions returns a naive list of auctions.
func (s *Service) ActiveAuctions(ctx context.Context, now time.Time) ([]Auction, error) {
    m, err := s.space.GetMap(ctx, "auctions")
    if err != nil {
        return nil, err
    }
    cctx, cancel := ihz.WithTimeout(ctx)
    defer cancel()
    kvs, err := m.GetEntrySet(cctx)
    if err != nil {
        return nil, err
    }
    out := make([]Auction, 0, len(kvs))
    for _, kv := range kvs {
        a, ok := kv.Value.(Auction)
        if !ok {
            continue
        }
        if now.After(a.StartsAt) && now.Before(a.EndsAt) {
            out = append(out, a)
        }
    }
    return out, nil
}


