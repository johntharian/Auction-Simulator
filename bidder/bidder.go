package bidder

import (
	"math/rand"
	"time"

	"auction-simulator/models"
)

// Bidder simulates a participant in auctions.
type Bidder struct {
	ID                  int
	weights             [20]float64 // internal valuation weights for each attribute
	responseProbability float64     // probability that this bidder responds (0.0–1.0)
	rng                 *rand.Rand
}

// NewBidders creates n bidders with randomized valuation profiles.
func NewBidders(n int, seed int64) []*Bidder {
	bidders := make([]*Bidder, n)
	src := rand.NewSource(seed)
	masterRng := rand.New(src)

	for i := 0; i < n; i++ {
		b := &Bidder{
			ID:                  i + 1,
			responseProbability: 0.5 + masterRng.Float64()*0.4, // 50%–90% chance of responding
			rng:                 rand.New(rand.NewSource(masterRng.Int63())),
		}
		for j := 0; j < 20; j++ {
			b.weights[j] = masterRng.Float64() * 2.0 // weight between 0 and 2
		}
		bidders[i] = b
	}
	return bidders
}

// PlaceBid evaluates the auction item and returns a bid.
// Returns (nil, false) if the bidder opts out of this auction.
func (b *Bidder) PlaceBid(item models.AuctionItem) (*models.Bid, bool) {
	// Decide whether to respond
	if b.rng.Float64() > b.responseProbability {
		return nil, false
	}

	// Simulate thinking time (5–50ms)
	delay := time.Duration(5+b.rng.Intn(46)) * time.Millisecond
	time.Sleep(delay)

	// Compute bid as weighted sum of attributes + noise
	var amount float64
	for i := 0; i < 20; i++ {
		amount += b.weights[i] * item.Attributes[i]
	}
	// Add some noise (±10%)
	noise := 1.0 + (b.rng.Float64()-0.5)*0.2
	amount *= noise

	if amount < 0 {
		amount = 0.01
	}

	return &models.Bid{
		BidderID:  b.ID,
		Amount:    amount,
		Timestamp: time.Now(),
	}, true
}
