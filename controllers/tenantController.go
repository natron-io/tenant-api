package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

func GetPods(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return c.Redirect("/login/github")
	}

	// create a map for each tenant with a list of pods with labels in it
	tenantPods, err := util.GetPodsByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantPods)
}

func GetNamespaces(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return c.Redirect("/login/github")
	}

	// create a map for each tenant with a list of namespaces with labels in it
	tenantNamespaces, err := util.GetNamespacesByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantNamespaces)
}

func GetServiceAccounts(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return c.Redirect("/login/github")
	}
	// create a map for each tenant with a map of namespaces with a list of service accounts with labels in it
	tenantServiceAccounts, err := util.GetServiceAccountsByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantServiceAccounts)
}

// get cpu request sum by tenant in millicores
func GetCPURequestsSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return c.Redirect("/login/github")
	}

	// create a map for each tenant with a added cpu requests
	tenantCPURequests, err := util.GetCPURequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantCPURequests)
}

// get memory request sum by tenant in bytes
func GetMemoryRequestsSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return c.Redirect("/login/github")
	}

	// create a map for each tenant with a added memory requests
	tenantMemoryRequests, err := util.GetMemoryRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantMemoryRequests)
}

// returns the sum in bytes of storagerequests by storageclass per tenant
func GetStorageRequestsSum(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return c.Redirect("/login/github")
	}

	// create a map for each tenant with a map of storage classes with calculated pvcs in it
	tenantPVCs, err := util.GetStorageRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantPVCs)
}
