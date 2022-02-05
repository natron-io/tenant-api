package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

// GetCPUQuota returns the CPU quota of a tenant by the label at the tenant config namespace
func GetCPUQuota(c *fiber.Ctx) error {

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

	cpuQuota, err := util.GetRessourceQuota(tenant, util.QUOTA_NAMESPACE_SUFFIX, util.QUOTA_CPU_LABEL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(cpuQuota)
}

// GetMemoryQuota returns the Memory quota of a tenant by the label at the tenant config namespace
func GetMemoryQuota(c *fiber.Ctx) error {

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

	memoryQuota, err := util.GetRessourceQuota(tenant, util.QUOTA_NAMESPACE_SUFFIX, util.QUOTA_MEMORY_LABEL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(memoryQuota)
}

// GetStorageQuota returns the Storage quota of a tenant by the label at the tenant config namespace
func GetStorageQuota(c *fiber.Ctx) error {

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

	storageQuota := make(map[string]float64)
	var err error

	for key, value := range util.QUOTA_STORAGE_LABEL {
		storageQuota[key], err = util.GetRessourceQuota(tenant, util.QUOTA_NAMESPACE_SUFFIX, value)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	}

	return c.JSON(storageQuota)
}
