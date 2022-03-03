package util

import (
	"fmt"
	"strings"

	"github.com/natron-io/tenant-api/database"
	"github.com/natron-io/tenant-api/models"
)

var (
	CPU_COST                 float64
	MEMORY_COST              float64
	STORAGE_COST             map[string]map[string]float64
	INGRESS_COST             float64
	INGRESS_COST_PER_DOMAIN  bool
	EXCLUDE_INGRESS_VCLUSTER bool
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
		return 0, fmt.Errorf("storage class %s not found", storageClass)
	}
	// with STORAGE_DISCOUNT_PERCENT
	return (STORAGE_COST[storageClass]["cost"] * float64(size) / (1024 * 1024 * 1024)) * (1 - STORAGE_DISCOUNT_PERCENT), nil
}

// GetIngressCost returns the cost of the provided Ingress
func GetIngressCostByDomain(hostnameStrings []string) float64 {

	var tenantIngressCostsPerDomainSum float64

	// define a set of domains
	domains := make(map[string]bool)

	for _, host := range hostnameStrings {
		// split string with .
		hostnameParts := strings.Split(host, ".")
		// get the 2 last parts of the hostname
		var domain string
		if len(hostnameParts) > 1 {
			domain = hostnameParts[len(hostnameParts)-2] + "." + hostnameParts[len(hostnameParts)-1]
		} else {
			domain = ""
		}
		// add the domain to the tenantIngressCosts map
		if domain != "" {
			domains[domain] = true
		} else {
			ErrorLogger.Printf("domain is not valid for hostname %s", host)
		}
	}

	// calculate the cost * count of domains
	for range domains {
		tenantIngressCostsPerDomainSum += INGRESS_COST * (1 - INGRESS_DISCOUNT_PERCENT)
	}

	return tenantIngressCostsPerDomainSum
}

// GetIngressCost returns the cost of the provided Ingress
func GetIngressCost(ingressCount int) float64 {
	return INGRESS_COST * float64(ingressCount) * (1 - INGRESS_DISCOUNT_PERCENT)
}

// SaveCostsDB saves all the costs of a tenant in the database
func SaveCostsToDB() error {
	var err error
	var tenants []models.Tenant
	var cpuCost models.CPUCost
	var memoryCost models.MemoryCost
	var ingressCost models.IngressCost

	db := database.DBConn

	// get all namespaces in the kubernetes cluster
	namespaces, err := GetNamespaces()
	if err != nil {
		return err
	}

	for _, namespace := range namespaces {
		tenant := models.Tenant{
			GitHubTeamSlug: namespace,
		}

		tenants = append(tenants, tenant)

	}
	tenantCPURequests, err := GetCPURequestsSumByTenant(namespaces)
	if err != nil {
		return err
	}

	tenantMemoryRequests, err := GetMemoryRequestsSumByTenant(namespaces)
	if err != nil {
		return err
	}

	tenantPVCs, err := GetStorageRequestsSumByTenant(namespaces)
	if err != nil {
		return err
	}

	tenantsIngressRequests, err := GetIngressRequestsSumByTenant(namespaces)
	if err != nil {
		return err
	}

	for _, tenant := range tenants {
		cpuCost = models.CPUCost{
			TenantId: tenant.Id,
			Value:    GetCPUCost(float64(tenantCPURequests[tenant.GitHubTeamSlug])),
		}

		err = db.Create(&cpuCost).Error
		if err != nil {
			return err
		}
		InfoLogger.Printf("CPU cost for tenant %s: %f", tenant.GitHubTeamSlug, cpuCost.Value)

		memoryCost = models.MemoryCost{
			TenantId: tenant.Id,
			Value:    GetMemoryCost(float64(tenantMemoryRequests[tenant.GitHubTeamSlug])),
		}

		err = db.Create(&memoryCost).Error
		if err != nil {
			return err
		}
		InfoLogger.Printf("Memory cost for tenant %s: %f", tenant.GitHubTeamSlug, memoryCost.Value)

		for storageClass, pvc := range tenantPVCs[tenant.GitHubTeamSlug] {
			storageClassCost, err := GetStorageCost(storageClass, float64(pvc))
			if err != nil {
				return err
			}
			storageCost := models.StorageCost{
				TenantId:     tenant.Id,
				StorageClass: storageClass,
				Value:        storageClassCost,
			}

			err = db.Create(&storageCost).Error
			if err != nil {
				return err
			}
			InfoLogger.Printf("Storage cost for tenant %s: %f %s", tenant.GitHubTeamSlug, storageCost.Value, storageClass)
		}

		if !INGRESS_COST_PER_DOMAIN {
			ingressCost = models.IngressCost{
				TenantId: tenant.Id,
				Value:    GetIngressCost(len(tenantsIngressRequests[tenant.GitHubTeamSlug])),
			}

			err = db.Create(&ingressCost).Error
			if err != nil {
				return err
			}
			InfoLogger.Printf("Ingress cost for tenant %s: %f", tenant.GitHubTeamSlug, ingressCost.Value)
		} else {
			ingressCost = models.IngressCost{
				TenantId: tenant.Id,
				Value:    GetIngressCostByDomain(tenantsIngressRequests[tenant.GitHubTeamSlug]),
			}

			err = db.Create(&ingressCost).Error
			if err != nil {
				return err
			}
			InfoLogger.Printf("Ingress cost for tenant %s: %f", tenant.GitHubTeamSlug, ingressCost.Value)
		}

	}

	return nil
}
