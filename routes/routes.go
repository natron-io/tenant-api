package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/controllers"
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

	// Tenants
	v1.Get("/tenants", controllers.GetTenants)

	// Every Tenant
	v1.Get("/pods", controllers.GetPods)
	v1.Get("/namespaces", controllers.GetNamespaces)
	v1.Get("/serviceAccounts", controllers.GetServiceAccounts)

	// Specific Tenant
	v1.Get(":tenant/pods", controllers.GetPods)
	v1.Get(":tenant/namespaces", controllers.GetNamespaces)
	v1.Get(":tenant/serviceAccounts", controllers.GetServiceAccounts)

	// Every Tenant
	requests := v1.Group("/requests")
	requests.Get("/cpu", controllers.GetCPURequestsSum)
	requests.Get("/memory", controllers.GetMemoryRequestsSum)
	requests.Get("/storage", controllers.GetStorageRequestsSum)
	requests.Get("/ingress", controllers.GetIngressRequestsSum)

	// Specific Tenant
	requests = v1.Group(":tenant/requests")
	requests.Get("/cpu", controllers.GetCPURequestsSum)
	requests.Get("/memory", controllers.GetMemoryRequestsSum)
	requests.Get("/storage", controllers.GetStorageRequestsSum)
	requests.Get("/ingress", controllers.GetIngressRequestsSum)

	// Every Tenant
	costs := v1.Group("/costs")
	costs.Get("/cpu", controllers.GetCPUCostSum)
	costs.Get("/memory", controllers.GetMemoryCostSum)
	costs.Get("/storage", controllers.GetStorageCostSum)
	costs.Get("/ingress", controllers.GetIngressCostSum)

	// Per tenant
	costs = v1.Group(":tenant/costs")
	costs.Get("/cpu", controllers.GetCPUCostSum)
	costs.Get("/memory", controllers.GetMemoryCostSum)
	costs.Get("/storage", controllers.GetStorageCostSum)
	costs.Get("/ingress", controllers.GetIngressCostSum)
}
