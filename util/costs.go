package util

var (
	CPU_COST     float64
	MEMORY_COST  float64
	STORAGE_COST []StorageClassCost
)

type StorageClassCost struct {
	StorageClass string
	Cost         float64
}

func GetCPUCost(millicores float64) float64 {
	// return per core
	return CPU_COST * float64(millicores) / 1000
}

func GetMemoryCost(memory float64) float64 {
	// return per GB
	return MEMORY_COST * float64(memory) / (1024 * 1024 * 1024)
}

func GetStorageCost(storageClass string, size float64) float64 {
	for _, storageClassCost := range STORAGE_COST {
		if storageClassCost.StorageClass == storageClass {
			// return per GB
			return storageClassCost.Cost * float64(size) / (1024 * 1024 * 1024)
		}
	}
	// return per GB
	return STORAGE_COST[0].Cost * float64(size) / (1024 * 1024 * 1024)
}
