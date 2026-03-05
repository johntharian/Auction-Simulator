package engine

import (
	"auction-simulator/bidder"
	"auction-simulator/models"
	"context"
	"math/rand"
	"sync"
	"time"
)

// GenerateItems creates n auction items with random attributes.
func GenerateItems(n int, seed int64) []models.AuctionItem {
	rng := rand.New(rand.NewSource(seed))
	items := make([]models.AuctionItem, n)
	for i := 0; i < n; i++ {
		item := models.AuctionItem{ID: i + 1}
		for j := 0; j < 20; j++ {
			item.Attributes[j] = rng.Float64() * 10.0 // attribute values 0–10
		}
		items[i] = item
	}
	return items
}

// RunAuction runs a single auction: sends the item to all bidders concurrently,
// collects bids within the timeout, and declares a winner.
func RunAuction(auctionID int, item models.AuctionItem, bidders []*bidder.Bidder, timeout time.Duration) models.AuctionResult {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	bidCh := make(chan models.Bid, len(bidders))
	var wg sync.WaitGroup

	// Launch all bidders concurrently
	for _, b := range bidders {
		wg.Add(1)
		go func(b *bidder.Bidder) {
			defer wg.Done()
			bid, ok := b.PlaceBid(item)
			if !ok {
				return // bidder opted out
			}
			select {
			case <-ctx.Done():
				return // timeout reached, discard late bid
			case bidCh <- *bid:
			}
		}(b)
	}

	// Close the bid channel once all bidders finish or context expires
	go func() {
		wg.Wait()
		close(bidCh)
	}()

	// Collect bids until the channel closes or timeout
	var allBids []models.Bid
	var winner *models.Bid

collectLoop:
	for {
		select {
		case bid, ok := <-bidCh:
			if !ok {
				break collectLoop
			}
			allBids = append(allBids, bid)
			if winner == nil || bid.Amount > winner.Amount {
				winnerCopy := bid
				winner = &winnerCopy
			}
		case <-ctx.Done():
			// Drain any remaining bids already in the channel
			for bid := range bidCh {
				allBids = append(allBids, bid)
				if winner == nil || bid.Amount > winner.Amount {
					winnerCopy := bid
					winner = &winnerCopy
				}
			}
			break collectLoop
		}
	}

	endTime := time.Now()
	status := "completed"
	if ctx.Err() == context.DeadlineExceeded {
		status = "timeout"
	}

	return models.AuctionResult{
		AuctionID:  auctionID,
		Item:       item,
		TimeoutMs:  int(timeout.Milliseconds()),
		TotalBids:  len(allBids),
		Winner:     winner,
		AllBids:    allBids,
		StartTime:  startTime,
		EndTime:    endTime,
		DurationMs: float64(endTime.Sub(startTime).Microseconds()) / 1000.0,
		Status:     status,
	}
}

// RunAllAuctions launches all auctions concurrently and measures overall duration.
func RunAllAuctions(items []models.AuctionItem, bidders []*bidder.Bidder, timeout time.Duration) ([]models.AuctionResult, time.Duration) {
	results := make([]models.AuctionResult, len(items))
	var wg sync.WaitGroup

	overallStart := time.Now()

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it models.AuctionItem) {
			defer wg.Done()
			results[idx] = RunAuction(idx+1, it, bidders, timeout)
		}(i, item)
	}

	wg.Wait()
	overallDuration := time.Since(overallStart)

	return results, overallDuration
}
