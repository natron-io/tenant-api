package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/controllers"
	"github.com/natron-io/tenant-api/util"
	"k8s.io/client-go/kubernetes"
)

// Routes - Define all routes
func Setup(app *fiber.App, clientset *kubernetes.Clientset) {
	// Auth
	app.Post("/login/github", controllers.FrontendGithubLogin)
	app.Get("/login/github", controllers.GithubLogin)
	app.Get("/login/github/callback", controllers.GithubCallback)

	// API
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Notifications
	if util.SLACK_TOKEN != "" {
		v1.Get("/notifications", controllers.GetNotifications)
	}

	// Tenants
	v1.Get("/tenants", controllers.GetTenants)

	// Specific Tenant
	v1.Get(":tenant/pods", controllers.GetPods)

	// Specific Tenant
	requests := v1.Group(":tenant/requests")
	requests.Get("/cpu", controllers.GetCPURequestsSum)
	requests.Get("/memory", controllers.GetMemoryRequestsSum)
	requests.Get("/storage", controllers.GetStorageRequestsSum)
	requests.Get("/ingress", controllers.GetIngressRequestsSum)

	// Per tenant
	costs := v1.Group(":tenant/costs")
	costs.Get("/cpu", controllers.GetCPUCostSum)
	costs.Get("/memory", controllers.GetMemoryCostSum)
	costs.Get("/storage", controllers.GetStorageCostSum)
	costs.Get("/ingress", controllers.GetIngressCostSum)

	// Quotas
	quotas := v1.Group(":tenant/quotas")
	quotas.Get("/cpu", controllers.GetCPUQuota)
	quotas.Get("/memory", controllers.GetMemoryQuota)
	quotas.Get("/storage", controllers.GetStorageQuota)
}
