package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/controllers"
	"k8s.io/client-go/kubernetes"
)

func Setup(app *fiber.App, clientset *kubernetes.Clientset) {
	// Auth
	app.Get("/login/github", controllers.GithubLogin)
	app.Get("/login/github/callback", controllers.GithubCallback)
	app.Get("/loggedin", func(c *fiber.Ctx) error {
		return controllers.LoggedIn(c, c.Get("githubData"))
	})

	// API
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/pods", controllers.GetPods)
	v1.Get("/namespaces", controllers.GetNamespaces)
	v1.Get("/serviceAccounts", controllers.GetServiceAccounts)
	v1.Get("/cpurequests", controllers.GetCPURequestsSum)
	v1.Get("/memoryrequests", controllers.GetMemoryRequestsSum)
	v1.Get("/storagerequests", controllers.GetStorageAllocationSum)
}
