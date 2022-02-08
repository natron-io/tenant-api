package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

// GetTenants returns all tenants by authentication
func GetTenants(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.JSON(tenants)
}

// GetPods returns all pods by authenticated users tenants
func GetPods(c *fiber.Ctx) error {

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

	var tenantPods map[string][]string
	var err error
	if tenant == "" {
		tenantPods, err = util.GetPodsByTenant(tenants)
		if err != nil {
			util.ErrorLogger.Printf("%s", err)
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		return c.JSON(tenantPods)
	} else {
		tenantPods, err = util.GetPodsByTenant([]string{tenant})
		if err != nil {
			util.ErrorLogger.Printf("%s", err)
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		return c.JSON(tenantPods[tenant])
	}
}

// GetCPURequestsSum returns the sum of all cpu requests by authenticated users tenants
func GetCPURequestsSum(c *fiber.Ctx) error {

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

	if tenant == "" {
		return c.JSON(tenantCPURequests)
	} else {
		return c.JSON(tenantCPURequests[tenant])
	}
}

// GetMemoryRequestsSum returns the sum of all memory requests by authenticated users tenants
func GetMemoryRequestsSum(c *fiber.Ctx) error {

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

	if tenant == "" {
		return c.JSON(tenantMemoryRequests)
	} else {
		return c.JSON(tenantMemoryRequests[tenant])
	}
}

// GetStorageRequestsSum returns the sum of all storage requests by authenticated users tenants
func GetStorageRequestsSum(c *fiber.Ctx) error {

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

	if tenant == "" {
		return c.JSON(tenantPVCs)
	} else {
		return c.JSON(tenantPVCs[tenant])
	}
}

// GetIngressRequestsSum returns the sum of all ingress requests by authenticated users tenants
func GetIngressRequestsSum(c *fiber.Ctx) error {

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
	tenantIngressRequests, err := util.GetIngressRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	if tenant == "" {
		return c.JSON(tenantIngressRequests)
	} else {
		return c.JSON(tenantIngressRequests[tenant])
	}
}
