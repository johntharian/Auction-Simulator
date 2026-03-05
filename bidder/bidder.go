package bidder

import (
	"math/rand"
	"time"

	"auction-simulator/models"
)

// Bidder simulates a participant in auctions. Each bidder has a unique set of
// valuation weights for item attributes, a response probability, and its own
// random number generator for deterministic behavior.
//
// Fields:
//   - ID: unique identifier for the bidder.
//   - weights: internal valuation weights (0–2) for each of the 20 item attributes.
//   - responseProbability: probability (0.5–0.9) that this bidder participates in a given auction.
//   - rng: per-bidder random number generator for independent randomness.
type Bidder struct {
	ID                  int
	weights             [20]float64 // internal valuation weights for each attribute
	responseProbability float64     // probability that this bidder responds (0.0–1.0)
	rng                 *rand.Rand
}

// NewBidders creates a slice of n bidders, each with randomized valuation weights
// and response probabilities derived from the provided seed. A master RNG is used
// to seed individual per-bidder RNGs, ensuring deterministic and reproducible
// bidder populations.
//
// Parameters:
//   - n: the number of bidders to create.
//   - seed: the seed for the master random number generator.
//
// Returns:
//   - []*Bidder: a slice of pointers to the newly created bidders.
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

// PlaceBid evaluates the given auction item and decides whether to place a bid.
// The bidder first randomly decides whether to participate based on its response
// probability. If participating, it simulates a thinking delay (5–50 ms), computes
// a bid amount as the weighted sum of the item's attributes (plus ±10% noise), and
// returns the resulting bid.
//
// Parameters:
//   - item: the auction item to evaluate and potentially bid on.
//
// Returns:
//   - *models.Bid: the bid placed by this bidder, or nil if the bidder opts out.
//   - bool: true if a bid was placed, false if the bidder opted out.
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
