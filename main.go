package main

import (
	"auction-simulator/bidder"
	"auction-simulator/engine"
	"auction-simulator/models"
	"auction-simulator/resource"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// --- CLI Flags ---
	vcpus := flag.Int("vcpus", 4, "Number of vCPUs (GOMAXPROCS)")
	memoryMB := flag.Int("memory-mb", 512, "Memory limit in MB (GOMEMLIMIT)")
	timeoutMs := flag.Int("timeout-ms", 200, "Per-auction timeout in milliseconds")
	outputDir := flag.String("output-dir", "output", "Directory to write auction result JSON files")
	numBidders := flag.Int("bidders", 100, "Number of bidders")
	numAuctions := flag.Int("auctions", 40, "Number of concurrent auctions")
	seed := flag.Int64("seed", 42, "Random seed for reproducibility")
	flag.Parse()

	// --- Apply Resource Constraints ---
	fmt.Println("========================================")
	fmt.Println("       AUCTION SIMULATOR")
	fmt.Println("========================================")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Bidders:       %d\n", *numBidders)
	fmt.Printf("  Auctions:      %d\n", *numAuctions)
	fmt.Printf("  Timeout:       %d ms\n", *timeoutMs)
	fmt.Printf("  vCPUs:         %d\n", *vcpus)
	fmt.Printf("  Memory Limit:  %d MB\n", *memoryMB)
	fmt.Printf("  Seed:          %d\n", *seed)
	fmt.Println("----------------------------------------")

	resource.ApplyConstraints(*vcpus, *memoryMB)
	fmt.Println("----------------------------------------")

	// --- Create Bidders ---
	bidders := bidder.NewBidders(*numBidders, *seed)
	fmt.Printf("Created %d bidders\n", len(bidders))

	// --- Generate Auction Items ---
	items := engine.GenerateItems(*numAuctions, *seed+1)
	fmt.Printf("Generated %d auction items (20 attributes each)\n", len(items))
	fmt.Println("----------------------------------------")

	// --- Run All Auctions Concurrently ---
	timeout := time.Duration(*timeoutMs) * time.Millisecond
	fmt.Printf("Starting %d auctions concurrently...\n\n", *numAuctions)

	results, totalDuration := engine.RunAllAuctions(items, bidders, timeout)

	// --- Report Results ---
	fmt.Println("========================================")
	fmt.Println("            RESULTS SUMMARY")
	fmt.Println("========================================")
	fmt.Printf("Total elapsed time: %v\n", totalDuration)
	fmt.Printf("Total elapsed (ms): %.2f\n", float64(totalDuration.Microseconds())/1000.0)
	fmt.Println("----------------------------------------")

	totalBids := 0
	for _, r := range results {
		totalBids += r.TotalBids
		winnerInfo := "No bids received"
		if r.Winner != nil {
			winnerInfo = fmt.Sprintf("Bidder #%d — $%.2f", r.Winner.BidderID, r.Winner.Amount)
		}
		fmt.Printf("  Auction %2d | %3d bids | %6.1f ms | %s | Winner: %s\n",
			r.AuctionID, r.TotalBids, r.DurationMs, r.Status, winnerInfo)
	}
	fmt.Println("----------------------------------------")
	fmt.Printf("Total bids across all auctions: %d\n", totalBids)

	// --- Resource Snapshot ---
	resInfo := resource.ReportResources()
	fmt.Println("----------------------------------------")
	fmt.Println("Resource Snapshot:")
	fmt.Printf("  GOMAXPROCS:     %d\n", resInfo.GOMAXPROCS)
	fmt.Printf("  Memory Limit:   %d MB\n", resInfo.MemoryLimitMB)
	fmt.Printf("  Current Alloc:  %.2f MB\n", resInfo.AllocMB)
	fmt.Printf("  System Memory:  %.2f MB\n", resInfo.SysMB)
	fmt.Printf("  Goroutines:     %d\n", resInfo.NumGoroutines)
	fmt.Println("========================================")

	// --- Write Output Files ---
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	for _, r := range results {
		writeAuctionResult(*outputDir, r)
	}
	fmt.Printf("\n✅ Wrote %d auction result files to %s/\n", len(results), *outputDir)
}

// writeAuctionResult writes a single auction result to a JSON file.
func writeAuctionResult(dir string, result models.AuctionResult) {
	filename := filepath.Join(dir, fmt.Sprintf("auction_%02d.json", result.AuctionID))
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling auction %d: %v\n", result.AuctionID, err)
		return
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", filename, err)
	}
}
