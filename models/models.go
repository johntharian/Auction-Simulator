package models

import "time"

// AuctionItem represents an object being auctioned, described by a fixed set of
// 20 numeric attributes. Bidders evaluate these attributes using their internal
// valuation weights to determine their bid amounts.
//
// Fields:
//   - ID: unique identifier for the auction item.
//   - Attributes: an array of 20 float64 values (range 0–10) representing the
//     item's characteristics.
type AuctionItem struct {
	ID         int        `json:"id"`
	Attributes [20]float64 `json:"attributes"`
}

// Bid represents a single bid placed by a bidder during an auction round.
//
// Fields:
//   - BidderID: the unique identifier of the bidder who placed this bid.
//   - Amount: the monetary amount of the bid.
//   - Timestamp: the wall-clock time at which the bid was created.
type Bid struct {
	BidderID  int       `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// AuctionResult captures the complete outcome of a single auction, including the
// item auctioned, all bids received, the winning bid, and timing/status metadata.
//
// Fields:
//   - AuctionID: unique identifier for this auction round.
//   - Item: the auction item that was bid on.
//   - TimeoutMs: the per-auction timeout in milliseconds.
//   - TotalBids: the total number of bids received.
//   - Winner: pointer to the winning (highest) bid, or nil if no bids were received.
//   - AllBids: a slice of all bids collected during the auction.
//   - StartTime: the wall-clock time the auction started.
//   - EndTime: the wall-clock time the auction ended.
//   - DurationMs: the auction duration in milliseconds (microsecond precision).
//   - Status: "completed" if all bidders responded, or "timeout" if the deadline was hit.
type AuctionResult struct {
	AuctionID      int          `json:"auction_id"`
	Item           AuctionItem  `json:"item"`
	TimeoutMs      int          `json:"timeout_ms"`
	TotalBids      int          `json:"total_bids_received"`
	Winner         *Bid         `json:"winner"`
	AllBids        []Bid        `json:"all_bids"`
	StartTime      time.Time    `json:"start_time"`
	EndTime        time.Time    `json:"end_time"`
	DurationMs     float64      `json:"duration_ms"`
	Status         string       `json:"status"`
}

// ResourceInfo captures a snapshot of the Go runtime's resource constraints and
// current memory/goroutine usage at a point in time.
//
// Fields:
//   - GOMAXPROCS: the current GOMAXPROCS setting (max OS threads for Go code).
//   - MemoryLimitMB: the soft memory limit (GOMEMLIMIT) in megabytes.
//   - AllocMB: the current heap allocation in megabytes.
//   - SysMB: the total system memory obtained from the OS in megabytes.
//   - NumGoroutines: the number of currently active goroutines.
type ResourceInfo struct {
	GOMAXPROCS     int    `json:"gomaxprocs"`
	MemoryLimitMB  int    `json:"memory_limit_mb"`
	AllocMB        float64 `json:"current_alloc_mb"`
	SysMB          float64 `json:"system_memory_mb"`
	NumGoroutines  int    `json:"num_goroutines"`
}
