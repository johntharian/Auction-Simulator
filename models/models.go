package models

import "time"

// AuctionItem represents an object being auctioned, described by 20 attributes.
type AuctionItem struct {
	ID         int        `json:"id"`
	Attributes [20]float64 `json:"attributes"`
}

// Bid represents a single bid placed by a bidder.
type Bid struct {
	BidderID  int       `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// AuctionResult captures the complete outcome of a single auction.
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

// ResourceInfo captures a snapshot of resource constraints and usage.
type ResourceInfo struct {
	GOMAXPROCS     int    `json:"gomaxprocs"`
	MemoryLimitMB  int    `json:"memory_limit_mb"`
	AllocMB        float64 `json:"current_alloc_mb"`
	SysMB          float64 `json:"system_memory_mb"`
	NumGoroutines  int    `json:"num_goroutines"`
}
