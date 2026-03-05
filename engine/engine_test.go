package engine

import (
	"auction-simulator/bidder"
	"auction-simulator/models"
	"testing"
	"time"
)

func TestGenerateItems(t *testing.T) {
	items := GenerateItems(40, 42)
	if len(items) != 40 {
		t.Fatalf("expected 40 items, got %d", len(items))
	}
	for _, item := range items {
		for j, attr := range item.Attributes {
			if attr < 0 || attr > 10 {
				t.Errorf("item %d attribute %d = %.2f, expected 0–10", item.ID, j, attr)
			}
		}
	}
}

func TestRunAuction_WinnerDeclared(t *testing.T) {
	bidders := bidder.NewBidders(10, 42)
	item := models.AuctionItem{ID: 1}
	for i := 0; i < 20; i++ {
		item.Attributes[i] = 5.0
	}

	result := RunAuction(1, item, bidders, 2*time.Second)

	if result.AuctionID != 1 {
		t.Errorf("expected auction ID 1, got %d", result.AuctionID)
	}
	if result.TotalBids == 0 {
		t.Error("expected at least some bids")
	}
	if result.Winner == nil {
		t.Error("expected a winner to be declared")
	}
	if result.Winner != nil {
		// Winner should have the highest bid
		for _, bid := range result.AllBids {
			if bid.Amount > result.Winner.Amount {
				t.Errorf("found bid %.2f higher than winner %.2f", bid.Amount, result.Winner.Amount)
			}
		}
	}
	t.Logf("Auction completed: %d bids, winner: bidder #%d ($%.2f), duration: %.1fms, status: %s",
		result.TotalBids, result.Winner.BidderID, result.Winner.Amount, result.DurationMs, result.Status)
}

func TestRunAuction_TimeoutEnforced(t *testing.T) {
	bidders := bidder.NewBidders(100, 42)
	item := models.AuctionItem{ID: 1}
	for i := 0; i < 20; i++ {
		item.Attributes[i] = 5.0
	}

	// Very short timeout — some bidders should be cut off
	result := RunAuction(1, item, bidders, 10*time.Millisecond)

	// With 10ms timeout and 5–50ms bidder delay, we expect fewer bids than bidders
	if result.TotalBids >= 100 {
		t.Errorf("expected fewer than 100 bids with 10ms timeout, got %d", result.TotalBids)
	}
	t.Logf("Timeout test: %d/%d bids received in 10ms, status: %s", result.TotalBids, 100, result.Status)
}

func TestRunAllAuctions_Concurrent(t *testing.T) {
	bidders := bidder.NewBidders(20, 42)
	items := GenerateItems(5, 99)

	results, elapsed := RunAllAuctions(items, bidders, 500*time.Millisecond)

	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}

	for i, r := range results {
		if r.AuctionID != i+1 {
			t.Errorf("result %d has auction ID %d", i, r.AuctionID)
		}
	}

	// Concurrent execution should be faster than sequential (5 × 500ms = 2.5s)
	if elapsed > 2*time.Second {
		t.Errorf("expected concurrent execution to be fast, took %v", elapsed)
	}
	t.Logf("5 auctions completed in %v", elapsed)
}
