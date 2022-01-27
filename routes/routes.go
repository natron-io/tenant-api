package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/controllers"
	"github.com/natron-io/tenant-api/util"
	"k8s.io/client-go/kubernetes"
)

func Setup(app *fiber.App, clientset *kubernetes.Clientset) {
	// Auth
	if util.FRONTENDAUTH_ENABLED {
		app.Post("/login/github", controllers.FrontendGithubLogin)
	} else {
		app.Get("/login/github", controllers.GithubLogin)
		app.Get("/logout", controllers.Logout)
		app.Get("/login/github/callback", controllers.GithubCallback)
	}

	// API
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/pods", controllers.GetPods)
	v1.Get("/namespaces", controllers.GetNamespaces)
	v1.Get("/serviceAccounts", controllers.GetServiceAccounts)

	requests := v1.Group("/requests")
	requests.Get("/cpu", controllers.GetCPURequestsSum)
	requests.Get("/memory", controllers.GetMemoryRequestsSum)
	requests.Get("/storage", controllers.GetStorageRequestsSum)
	requests.Get("/ingress", controllers.GetIngressRequestsSum)

	costs := v1.Group("/costs")
	costs.Get("/cpu", controllers.GetCPUCostSum)
	costs.Get("/memory", controllers.GetMemoryCostSum)
	costs.Get("/storage", controllers.GetStorageCostSum)
	costs.Get("/ingress", controllers.GetIngressCostSum)
}
