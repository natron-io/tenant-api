package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

func GetTenants(c *fiber.Ctx) error {

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

	return c.JSON(tenants)
}

func GetPods(c *fiber.Ctx) error {

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

	// create a map for each tenant with a list of pods with labels in it
	tenantPods, err := util.GetPodsByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantPods)
}

func GetTenantPods(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a list of pods with labels in it
	tenantPods, err := util.GetPodsByTenant([]string{tenant})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantPods[tenant])
}

func GetNamespaces(c *fiber.Ctx) error {

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

	// create a map for each tenant with a list of namespaces with labels in it
	tenantNamespaces, err := util.GetNamespacesByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantNamespaces)
}

func GetTenantNamespaces(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a list of namespaces with labels in it
	tenantNamespaces, err := util.GetNamespacesByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantNamespaces[tenant])
}

func GetServiceAccounts(c *fiber.Ctx) error {

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
	// create a map for each tenant with a map of namespaces with a list of service accounts with labels in it
	tenantServiceAccounts, err := util.GetServiceAccountsByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantServiceAccounts)
}

func GetTenantServiceAccounts(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a map of namespaces with a list of service accounts with labels in it
	tenantServiceAccounts, err := util.GetServiceAccountsByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantServiceAccounts[tenant])
}

// get cpu request sum by tenant in millicores
func GetCPURequestsSum(c *fiber.Ctx) error {

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

	return c.JSON(tenantCPURequests)
}

func GetTenantCPURequestsSum(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a added cpu requests
	tenantCPURequests, err := util.GetCPURequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantCPURequests[tenant])
}

// get memory request sum by tenant in bytes
func GetMemoryRequestsSum(c *fiber.Ctx) error {

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

	return c.JSON(tenantMemoryRequests)
}

func GetTenantMemoryRequestsSum(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a added memory requests
	tenantMemoryRequests, err := util.GetMemoryRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantMemoryRequests[tenant])
}

// returns the sum in bytes of storagerequests by storageclass per tenant
func GetStorageRequestsSum(c *fiber.Ctx) error {

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

	return c.JSON(tenantPVCs)
}

func GetTenantStorageRequestsSum(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a map of storage classes with calculated pvcs in it
	tenantPVCs, err := util.GetStorageRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantPVCs[tenant])
}

func GetIngressRequestsSum(c *fiber.Ctx) error {

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
	tenantIngressRequests, err := util.GetIngressRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantIngressRequests)
}

func GetTenantIngressRequestsSum(c *fiber.Ctx) error {
	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenant := c.Params("tenant")
	tenants := CheckAuth(c)
	if tenants == nil {
		c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		if !util.FRONTENDAUTH_ENABLED {
			return c.Redirect("/login/github")
		}
		// check if tenant is in the list of tenants
		if !util.Contains(tenant, tenants) {
			return c.Status(403).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tenant is not specified",
		})
	}

	// create a map for each tenant with a map of storage classes with calculated pvcs in it
	tenantIngressRequests, err := util.GetIngressRequestsSumByTenant(tenants)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(tenantIngressRequests[tenant])
}
