package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
	v1 "k8s.io/api/core/v1"
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

	quota, err := util.GetRessourceQuota(tenant)
	if err != nil {
		util.ErrorLogger.Printf("%s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	cpuQuota := quota.Spec.Hard.Cpu().MilliValue()

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

	quota, err := util.GetRessourceQuota(tenant)
	if err != nil {
		util.ErrorLogger.Printf("%s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	memoryQuota := quota.Spec.Hard.Memory().Value()

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

	quota, err := util.GetRessourceQuota(tenant)

	if err != nil {
		util.ErrorLogger.Printf("%s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	storageClasses, err := util.GetStorageClassesInCluster()
	if err != nil {
		util.ErrorLogger.Printf("%s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	storageQuota := quota.Spec.Hard

	// get first element of storageQuota
	storageQuotaMap := make(map[v1.ResourceName]int64)
	for key, value := range storageQuota {
		storageQuotaMap[key] = value.Value()
	}

	storageQuotaParsed := make(map[string]int64)
	for _, storageClass := range storageClasses {
		util.InfoLogger.Printf("%s", storageClass)
		storageQuotaString := v1.ResourceName(storageClass + ".storageclass.storage.k8s.io/requests.storage")
		if _, ok := storageQuotaMap[storageQuotaString]; ok {
			// convert to bytes
			storageQuotaParsed[storageClass] = storageQuotaMap[storageQuotaString] * 1024 * 1024
		} else {
			storageQuotaParsed[storageClass] = 0
		}
	}

	// check if storageClass string is in storageQuota
	return c.JSON(storageQuotaParsed)
}
