package auction

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/suite"
    "golang-sba-hazelcast/internal/platform/testdouble"
)

type serviceSuite struct {
    suite.Suite
    svc *Service
    ctx context.Context
}

func TestServiceSuite(t *testing.T) {
    suite.Run(t, new(serviceSuite))
}

func (s *serviceSuite) SetupTest() {
    f := testdouble.NewFakeSpace()
    s.svc = NewService(f)
    s.ctx = context.Background()
}

func (s *serviceSuite) TestCreateAuction_Validation() {
    err := s.svc.CreateAuction(s.ctx, Auction{})
    s.Error(err)
}

func (s *serviceSuite) TestCreateAndListActive() {
    now := time.Now()
    a := Auction{
        ID:          "a1",
        Title:       "Test",
        Description: "desc",
        StartsAt:    now.Add(-time.Minute),
        EndsAt:      now.Add(time.Minute),
        SellerID:    "u1",
    }
    s.NoError(s.svc.CreateAuction(s.ctx, a))
    list, err := s.svc.ActiveAuctions(s.ctx, now)
    s.NoError(err)
    s.Len(list, 1)
}

func (s *serviceSuite) TestPlaceBid() {
    b := Bid{AuctionID: "a1", UserID: "u2", Amount: 100, PlacedAt: time.Now()}
    s.NoError(s.svc.PlaceBid(s.ctx, b))
}


