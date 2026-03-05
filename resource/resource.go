package resource

import (
	"auction-simulator/models"
	"fmt"
	"runtime"
	"runtime/debug"
)

// ApplyConstraints configures runtime resource constraints for the simulation.
// It sets GOMAXPROCS to cap CPU parallelism and GOMEMLIMIT to impose a soft memory
// ceiling on the Go garbage collector. Values of zero or negative are ignored.
//
// Parameters:
//   - vcpus: the desired GOMAXPROCS value (number of OS threads for Go code).
//   - memoryMB: the soft memory limit in megabytes applied via debug.SetMemoryLimit.
func ApplyConstraints(vcpus, memoryMB int) {
	if vcpus > 0 {
		runtime.GOMAXPROCS(vcpus)
		fmt.Printf("[resource] GOMAXPROCS set to %d\n", vcpus)
	}
	if memoryMB > 0 {
		limit := int64(memoryMB) * 1024 * 1024
		debug.SetMemoryLimit(limit)
		fmt.Printf("[resource] Memory limit set to %d MB\n", memoryMB)
	}
}

// ReportResources reads the current Go runtime state and returns a snapshot of
// resource configuration (GOMAXPROCS, memory limit) and live usage metrics
// (heap allocation, system memory, goroutine count).
//
// Returns:
//   - models.ResourceInfo: a struct containing the current resource configuration
//     and memory/goroutine usage statistics.
func ReportResources() models.ResourceInfo {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return models.ResourceInfo{
		GOMAXPROCS:    runtime.GOMAXPROCS(0),
		MemoryLimitMB: int(debug.SetMemoryLimit(-1) / (1024 * 1024)), // read current limit
		AllocMB:       float64(memStats.Alloc) / (1024 * 1024),
		SysMB:         float64(memStats.Sys) / (1024 * 1024),
		NumGoroutines: runtime.NumGoroutine(),
	}
}
