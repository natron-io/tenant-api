package util

import (
	"fmt"
)

var (
	CPU_COST                 float64
	MEMORY_COST              float64
	STORAGE_COST             map[string]map[string]float64
	INGRESS_COST             float64
	CPU_DISCOUNT_PERCENT     float64
	MEMORY_DISCOUNT_PERCENT  float64
	STORAGE_DISCOUNT_PERCENT float64
	INGRESS_DISCOUNT_PERCENT float64
)

// GetCPUCost returns the cost of the provided MiliCPU
func GetCPUCost(millicores float64) float64 {
	// return per core
	return (CPU_COST * float64(millicores) / 1000) * (1 - CPU_DISCOUNT_PERCENT)
}

// GetMemoryCost returns the cost of the provided Memory
func GetMemoryCost(memory float64) float64 {
	// return per GB
	return (MEMORY_COST * float64(memory) / (1024 * 1024 * 1024)) * (1 - MEMORY_DISCOUNT_PERCENT)
}

// GetStorageCost returns the cost of the provided Storage of the StorageClass
func GetStorageCost(storageClass string, size float64) (float64, error) {
	// return per GB
	if STORAGE_COST[storageClass] == nil {
		Status = "error"
		return 0, fmt.Errorf("storage class %s not found", storageClass)
	}
	// with STORAGE_DISCOUNT_PERCENT
	return (STORAGE_COST[storageClass]["cost"] * float64(size) / (1024 * 1024 * 1024)) * (1 - STORAGE_DISCOUNT_PERCENT), nil
}

// GetIngressCost returns the cost of the provided Ingress
func GetIngressCost(ingress int) float64 {
	return (INGRESS_COST * float64(ingress)) * (1 - INGRESS_DISCOUNT_PERCENT)
}
