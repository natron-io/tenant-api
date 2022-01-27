package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

func GetCPUCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
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

	return c.JSON(tenantCPUCosts)
}

func GetMemoryCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
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

	return c.JSON(tenantMemoryCosts)
}

func GetStorageCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
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

	return c.JSON(tenantStorageCosts)
}

func GetIngressCostSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
	}

	// create a map for each tenant with a added ingress requests
	tenantIngressRequests, err := util.GetIngressRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// create a map for each tenant with a added ingress costs only if cost is not 0
	tenantIngressCosts := make(map[string]float64)
	for _, tenant := range tenants {
		if tenantIngressRequests[tenant] != 0 {
			tenantIngressCosts[tenant] = util.GetIngressCost(tenantIngressRequests[tenant])
		}
	}

	return c.JSON(tenantIngressCosts)
}
