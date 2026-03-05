package bidder

import (
	"auction-simulator/models"
	"testing"
)

func TestNewBidders(t *testing.T) {
	bidders := NewBidders(100, 42)
	if len(bidders) != 100 {
		t.Fatalf("expected 100 bidders, got %d", len(bidders))
	}
	// Check IDs are 1-indexed
	for i, b := range bidders {
		if b.ID != i+1 {
			t.Errorf("bidder %d has ID %d, expected %d", i, b.ID, i+1)
		}
	}
	// Verify response probability is within range
	for _, b := range bidders {
		if b.responseProbability < 0.5 || b.responseProbability > 0.9 {
			t.Errorf("bidder %d has responseProbability %.3f, expected 0.5–0.9", b.ID, b.responseProbability)
		}
	}
}

func TestPlaceBid_Deterministic(t *testing.T) {
	_ = NewBidders(1, 123) // verify no panic
	item := models.AuctionItem{ID: 1}
	for i := 0; i < 20; i++ {
		item.Attributes[i] = float64(i+1) * 0.5
	}

	// Multiple calls to test that bids are generated (at least some should respond)
	responded := 0
	for i := 0; i < 20; i++ {
		reconfigured := NewBidders(1, int64(1000+i))
		bid, ok := reconfigured[0].PlaceBid(item)
		if ok {
			responded++
			if bid.BidderID != 1 {
				t.Errorf("expected bidder ID 1, got %d", bid.BidderID)
			}
			if bid.Amount <= 0 {
				t.Errorf("expected positive bid amount, got %f", bid.Amount)
			}
		}
	}
	if responded == 0 {
		t.Error("expected at least some bids across 20 trials, got none")
	}
}

func TestPlaceBid_OptOut(t *testing.T) {
	// With enough trials, at least some bidders should opt out
	bidders := NewBidders(100, 99)
	item := models.AuctionItem{ID: 1}
	for i := 0; i < 20; i++ {
		item.Attributes[i] = 1.0
	}

	optedOut := 0
	for _, b := range bidders {
		_, ok := b.PlaceBid(item)
		if !ok {
			optedOut++
		}
	}
	if optedOut == 0 {
		t.Error("expected at least some bidders to opt out")
	}
	t.Logf("%d out of %d bidders opted out", optedOut, len(bidders))
}
