/*
Copyright 2022 Jan Lauber

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/natron-io/tenant-api/routes"
	"github.com/natron-io/tenant-api/util"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {
	util.InitLoggers()
	util.Status = "Running"

	// load util config envs
	if err := util.LoadEnv(); err != nil {
		util.ErrorLogger.Println("Error loading env variables")
		os.Exit(1)
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		util.ErrorLogger.Printf("Error creating in-cluster config: %v", err)
		os.Exit(1)
	}
	// creates the clientset
	util.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		util.ErrorLogger.Printf("Error creating clientset: %v", err)
		os.Exit(1)
	}
}

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET",
		AllowCredentials: true,
	}))

	app.Static("/styles", "./static/styles")
	app.Static("/images", "./static/images")

	app.Get("/", func(c *fiber.Ctx) error {
		// set header to html
		c.Set("Content-Type", "text/html") //TODO render css
		return c.Render("index", fiber.Map{
			"title":  "Tenant API",
			"status": util.GetStatus(),
		})
	})

	routes.Setup(app, util.Clientset)

	util.InfoLogger.Println("Tenant API is running on port 8000")

	app.Listen(":8000")
}
