package util

import (
	"fmt"
)

var (
	CPU_COST     float64
	MEMORY_COST  float64
	STORAGE_COST map[string]float64
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
	InfoLogger.Printf("GetStorageCost: %s %f", storageClass, size)
	// print STORAGE_COST
	for key, value := range STORAGE_COST {
		fmt.Printf("%s: %f\n", key, value)
	}
	// return per GB
	if cost, ok := STORAGE_COST[storageClass]; ok {
		return cost * float64(size) / (1024 * 1024 * 1024), nil
	}
	return 0, fmt.Errorf("storage class %s not found", storageClass)
}
