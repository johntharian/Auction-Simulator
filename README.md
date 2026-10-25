# Auction Simulator

A concurrent auction simulator written in **Go** that runs multiple auctions simultaneously, collects bids from simulated bidders, and measures execution times while adhering to configurable resource constraints.

## Architecture

```
auction-simulator/
├── main.go              # CLI entry point & orchestrator
├── models/models.go     # Core data types (AuctionItem, Bid, AuctionResult)
├── bidder/bidder.go      # Bidder simulation (valuation, latency, opt-out)
├── engine/engine.go      # Auction engine (concurrent execution, timeout, winner)
├── resource/resource.go  # Resource standardization (vCPU, memory)
└── output/               # JSON result files (one per auction)
```

## How It Works

1. **100 bidders** are created with random valuation profiles (weighted preferences across 20 attributes) and individual response probabilities (50%–90%).
2. **40 auction items** are generated, each described by 20 random attributes.
3. All **40 auctions run concurrently**. Within each auction, all 100 bidders are invoked in parallel via goroutines.
4. Each auction enforces a **timeout** using `context.WithTimeout`. Bids received after the deadline are discarded.
5. The **highest bid** wins each auction. Results are written as JSON files.
6. **Total elapsed time** (first auction start → last auction completion) is measured and reported.

## Resource Standardization

The simulator standardizes compute resources via two Go runtime controls:

| Flag           | Go Runtime API              | Purpose                              |
|----------------|-----------------------------|--------------------------------------|
| `--vcpus`      | `runtime.GOMAXPROCS(n)`     | Limits OS threads executing Go code  |
| `--memory-mb`  | `debug.SetMemoryLimit(n)`   | Sets soft memory limit (GOMEMLIMIT)  |

This ensures reproducible benchmarks across different hardware.

## Usage

### Build & Run

```bash
go build -o auction-simulator .
```

Run with default settings (100 bidders, 40 auctions, 200ms timeout):

```bash
./auction-simulator
```

Run with more bidders and a longer timeout:

```bash
./auction-simulator --bidders=200 --timeout-ms=500
```

Stress test with higher concurrency and constrained resources:

```bash
./auction-simulator --bidders=500 --auctions=100 --vcpus=2 --memory-mb=256 --timeout-ms=100
```

Quick test run with fewer auctions:

```bash
./auction-simulator --bidders=10 --auctions=5 --timeout-ms=1000
```

Run with a custom seed and output directory:

```bash
./auction-simulator --seed=123 --output-dir=results
```

### CLI Flags

| Flag            | Default | Description                            |
|-----------------|---------|----------------------------------------|
| `--vcpus`       | 4       | Number of vCPUs (GOMAXPROCS)           |
| `--memory-mb`   | 512     | Memory limit in MB                     |
| `--timeout-ms`  | 200     | Per-auction timeout in milliseconds    |
| `--output-dir`  | output  | Directory for auction result JSONs     |
| `--bidders`     | 100     | Number of bidders                      |
| `--auctions`    | 40      | Number of concurrent auctions          |
| `--seed`        | 42      | Random seed for reproducibility        |

### Run Tests

```bash
go test ./... -v
```

## Sample Output

Each auction produces a JSON file (`output/auction_01.json`, etc.):

```json
{
  "auction_id": 1,
  "item": {
    "id": 1,
    "attributes": [7.22, 3.15, 0.89, ...]
  },
  "timeout_ms": 200,
  "total_bids_received": 73,
  "winner": {
    "bidder_id": 42,
    "amount": 15.87,
    "timestamp": "2026-03-03T14:50:00.123Z"
  },
  "all_bids": [ ... ],
  "start_time": "2026-03-03T14:50:00.001Z",
  "end_time": "2026-03-03T14:50:00.201Z",
  "duration_ms": 48.2,
  "status": "completed"
}
```
