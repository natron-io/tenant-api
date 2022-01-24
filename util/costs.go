package util

import (
	"fmt"
)

var (
	CPU_COST     float64
	MEMORY_COST  float64
	STORAGE_COST map[string]map[string]float64
)

func GetCPUCost(millicores float64) float64 {
	// return per core
	return CPU_COST * float64(millicores) / 1000
}

func GetMemoryCost(memory float64) float64 {
	// return per GB
	return MEMORY_COST * float64(memory) / (1024 * 1024 * 1024)
}

func GetStorageCost(storageClass string, size float64) (float64, error) {
	// return per GB
	if STORAGE_COST[storageClass] == nil {
		Status = "error"
		return 0, fmt.Errorf("storage class %s not found", storageClass)
	}
	return STORAGE_COST[storageClass]["cost"] * float64(size) / (1024 * 1024 * 1024), nil
}
