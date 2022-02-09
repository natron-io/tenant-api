package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

// GetCPUCostSum returns the cpu cost sum per tenant
func GetCPUCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())
	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if len(tenants) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	if tenant != "" && !util.Contains(tenant, tenants) {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	// create a map for each tenant with a added cpu requests
	tenantCPURequests, err := util.GetCPURequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// create a map for each tenant with a added cpu costs only if cost is not 0
	tenantCPUCosts := make(map[string]float64)
	for _, tenant := range tenants {
		if tenantCPURequests[tenant] != 0 {
			tenantCPUCosts[tenant] = util.GetCPUCost(float64(tenantCPURequests[tenant]))
		}
	}

	if tenant == "" {
		return c.JSON(tenantCPUCosts)
	} else {
		return c.JSON(tenantCPUCosts[tenant])
	}
}

// GetMemoryCostSum returns the memory cost sum per tenant
func GetMemoryCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())
	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if len(tenants) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	if tenant != "" && !util.Contains(tenant, tenants) {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	// create a map for each tenant with a added memory requests
	tenantMemoryRequests, err := util.GetMemoryRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// create a map for each tenant with a added memory costs only if cost is not 0
	tenantMemoryCosts := make(map[string]float64)
	for _, tenant := range tenants {
		if tenantMemoryRequests[tenant] != 0 {
			tenantMemoryCosts[tenant] = util.GetMemoryCost(float64(tenantMemoryRequests[tenant]))
		}
	}

	if tenant == "" {
		return c.JSON(tenantMemoryCosts)
	} else {
		return c.JSON(tenantMemoryCosts[tenant])
	}
}

// GetStorageCostSum returns the storage cost sum per tenant
func GetStorageCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())
	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if len(tenants) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	if tenant != "" && !util.Contains(tenant, tenants) {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	// create a map for each tenant with a map of storage classes with calculated pvcs in it
	tenantPVCs, err := util.GetStorageRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// create a map for each tenant with each storage class with a cost if it is not 0 and add it to the tenant map
	tenantStorageCosts := make(map[string]map[string]float64)
	for _, tenant := range tenants {
		tenantStorageCosts[tenant] = make(map[string]float64)
		for storageClass, pvcs := range tenantPVCs[tenant] {
			if pvcs != 0 {
				tenantStorageCosts[tenant][storageClass], err = util.GetStorageCost(storageClass, float64(pvcs))
				if err != nil {
					util.ErrorLogger.Printf("%s", err)
					return c.Status(500).JSON(fiber.Map{
						"message": "Internal Server Error",
					})
				}
			}
		}
	}

	// remove tenants with no storage costs
	for tenant, storageCosts := range tenantStorageCosts {
		if len(storageCosts) == 0 {
			delete(tenantStorageCosts, tenant)
		}
	}

	if tenant == "" {
		return c.JSON(tenantStorageCosts)
	} else {
		return c.JSON(tenantStorageCosts[tenant])
	}
}

// GetIngressCostSum returns the ingress cost sum per tenant
func GetIngressCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())
	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if len(tenants) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	if tenant != "" && !util.Contains(tenant, tenants) {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	// create a map for each tenant with a added ingress requests
	tenantsIngressRequests, err := util.GetIngressRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// create a map for each tenant with a added ingress costs only if cost is not 0

	if !util.INGRESS_COST_PER_DOMAIN {
		tenantsIngressCosts := make(map[string]float64)
		for _, tenant := range tenants {
			if tenantsIngressRequests[tenant] != nil {
				tenantsIngressCosts[tenant] = util.GetIngressCost(len(tenantsIngressRequests[tenant]))
			}
		}

		if tenant == "" {
			return c.JSON(tenantsIngressCosts)
		} else {
			return c.JSON(tenantsIngressCosts[tenant])
		}
	}

	tenantsIngressCostsPerDomain := make(map[string]float64)
	for _, tenant := range tenants {
		if tenantsIngressRequests[tenant] != nil {
			tenantsIngressCostsPerDomain[tenant] = util.GetIngressCostByDomain(tenantsIngressRequests[tenant])
		}
	}

	if tenant == "" {
		return c.JSON(tenantsIngressCostsPerDomain)
	} else {
		return c.JSON(tenantsIngressCostsPerDomain[tenant])
	}
}
