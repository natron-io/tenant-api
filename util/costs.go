package util

import (
	"fmt"
	"strings"
	"time"

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

type costSum struct {
	Sum   float64
	Count int32
}

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

	// add or update the tenants in the database
	for _, tenant := range tenants {
		// get the tenant from the database
		err = db.Where("git_hub_team_slug = ?", tenant.GitHubTeamSlug).First(&tenant).Error
		if err != nil {
			// if the tenant is not in the database, add it
			err = db.Create(&tenant).Error
			if err != nil {
				return err
			}
		}
	}

	// get all tenants from the database
	err = db.Find(&tenants).Error
	if err != nil {
		return err
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

		memoryCost = models.MemoryCost{
			TenantId: tenant.Id,
			Value:    GetMemoryCost(float64(tenantMemoryRequests[tenant.GitHubTeamSlug])),
		}

		err = db.Create(&memoryCost).Error
		if err != nil {
			return err
		}

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

//
func GetCPUCostsByMonthAndTenant(tenant models.Tenant) (models.CPUCost, error) {

	db := database.DBConn

	// get all costs for the tenant for the current month
	var costs []models.CPUCost
	err := db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenant.Id, time.Now().AddDate(0, -1, 0), time.Now()).Find(&costs).Error
	if err != nil {
		return models.CPUCost{}, err
	}

	// calculate average cost for the month
	var sum float64
	var count int
	var averageCost float64
	for _, cost := range costs {
		sum += cost.Value
		count++
	}
	if count == 0 {
		averageCost = 0
	} else {
		averageCost = sum / float64(count)
	}

	cpuCosts := models.CPUCost{
		TenantId: tenant.Id,
		Value:    averageCost,
	}

	return cpuCosts, nil
}

func GetMemoryCostsByMonthAndTenant(tenant models.Tenant) (models.MemoryCost, error) {

	db := database.DBConn

	// get every cost for the tenant for the current month day 0 until now
	var costs []models.MemoryCost
	err := db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenant.Id, time.Now().AddDate(0, -1, 0), time.Now()).Find(&costs).Error
	if err != nil {
		return models.MemoryCost{}, err
	}

	// calculate average cost for the month
	var sum float64
	var count int
	var averageCost float64

	for _, cost := range costs {
		sum += cost.Value
		count++
	}
	if count == 0 {
		averageCost = 0
	} else {
		averageCost = sum / float64(count)
	}

	memoryCosts := models.MemoryCost{
		TenantId: tenant.Id,
		Value:    averageCost,
	}

	return memoryCosts, nil
}

func GetStorageCostsByMonthAndTenant(tenant models.Tenant) ([]models.StorageCost, error) {

	storageCostsTemp := make(map[string]costSum)
	storageCosts := make([]models.StorageCost, 0)

	db := database.DBConn

	// get all costs for the tenant for the current month
	var costs []models.StorageCost

	err := db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenant.Id, time.Now().AddDate(0, -1, 0), time.Now()).Find(&costs).Error
	if err != nil {
		return nil, err
	}

	// calculate average cost for the month for each storage class
	var averageCost float64

	for _, cost := range costs {
		storageCostsTemp[cost.StorageClass] = costSum{
			Sum:   storageCostsTemp[cost.StorageClass].Sum + cost.Value,
			Count: storageCostsTemp[cost.StorageClass].Count + 1,
		}
	}

	for storageClass, cost := range storageCostsTemp {
		averageCost = cost.Sum / float64(cost.Count)
		storageCosts = append(storageCosts, models.StorageCost{
			TenantId:     tenant.Id,
			StorageClass: storageClass,
			Value:        averageCost,
		})
	}

	return storageCosts, nil
}

func GetIngressCostsByMonthAndTenant(tenant models.Tenant) (models.IngressCost, error) {

	db := database.DBConn

	// get all costs for the tenant for the current month
	var costs []models.IngressCost
	err := db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenant.Id, time.Now().AddDate(0, -1, 0), time.Now()).Find(&costs).Error
	if err != nil {
		return models.IngressCost{}, err
	}

	// calculate average cost for the month
	var sum float64
	var count int
	var averageCost float64

	for _, cost := range costs {
		sum += cost.Value
		count++
	}

	if count == 0 {
		averageCost = 0
	} else {
		averageCost = sum / float64(count)
	}

	ingressCost := models.IngressCost{
		TenantId: tenant.Id,
		Value:    averageCost,
	}

	return ingressCost, nil

}

func CreateMonthlyCostReport() error {
	var err error
	var tenants []models.Tenant
	var monthlyCosts []models.MonthlyCost

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

	// add or update the tenants in the database
	for _, tenant := range tenants {
		// get the tenant from the database
		err = db.Where("git_hub_team_slug = ?", tenant.GitHubTeamSlug).First(&tenant).Error
		if err != nil {
			// if the tenant is not in the database, add it
			err = db.Create(&tenant).Error
			if err != nil {
				return err
			}
		}
	}

	// get all tenants from the database
	err = db.Find(&tenants).Error
	if err != nil {
		return err
	}

	for _, tenant := range tenants {
		cpuCosts, err := GetCPUCostsByMonthAndTenant(tenant)
		if err != nil {
			return err
		}

		memoryCosts, err := GetMemoryCostsByMonthAndTenant(tenant)
		if err != nil {
			return err
		}

		storageCosts, err := GetStorageCostsByMonthAndTenant(tenant)
		if err != nil {
			return err
		}
		var storageClassCostSum float64
		for _, storageClass := range storageCosts {
			storageClassCostSum = storageClassCostSum + storageClass.Value
		}

		ingressCosts, err := GetIngressCostsByMonthAndTenant(tenant)
		if err != nil {
			return err
		}

		// get last month and year
		lastMonth := time.Now().AddDate(0, -1, 0)
		lastYear := lastMonth.Year()

		// convert month and year to int
		monthInt := int32(lastMonth.Month())
		yearInt := int32(lastYear)

		monthlyCosts = append(monthlyCosts, models.MonthlyCost{
			TenantId:    tenant.Id,
			CPUCost:     cpuCosts.Value,
			MemoryCost:  memoryCosts.Value,
			StorageCost: storageClassCostSum,
			IngressCost: ingressCosts.Value,
			Month:       monthInt,
			Year:        yearInt,
			TotalCost:   cpuCosts.Value + memoryCosts.Value + storageClassCostSum + ingressCosts.Value,
		})
	}

	// add or update the monthly costs in the database
	for _, monthlyCost := range monthlyCosts {
		// get the monthly cost from the database
		err = db.Where("tenant_id = ? AND month = ? AND year = ?", monthlyCost.TenantId, monthlyCost.Month, monthlyCost.Year).First(&monthlyCost).Error
		if err != nil {
			// if the monthly cost is not in the database, add it
			err = db.Create(&monthlyCost).Error
			if err != nil {
				return err
			}

		} else {
			// if the monthly cost is in the database, update it
			err = db.Save(&monthlyCost).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
