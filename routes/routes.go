package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/controllers"
	"k8s.io/client-go/kubernetes"
)

func Setup(app *fiber.App, clientset *kubernetes.Clientset) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/pods", controllers.GetPods)
}
