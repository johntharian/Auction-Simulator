package resource

import (
	"auction-simulator/models"
	"fmt"
	"runtime"
	"runtime/debug"
)

// ApplyConstraints standardizes the available vCPU and RAM for the simulation.
//   - vcpus: sets GOMAXPROCS to limit the number of OS threads executing Go code.
//   - memoryMB: sets the Go runtime soft memory limit (GOMEMLIMIT) in megabytes.
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

// ReportResources returns a snapshot of current resource configuration and usage.
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
