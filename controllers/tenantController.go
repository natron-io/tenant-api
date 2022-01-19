package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(c *fiber.Ctx) error {
	// Get all pods in all namespaces
	pods, err := util.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	util.InfoLogger.Println("/api/v1/pods hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}
	// Return pods as JSON
	return c.JSON(pods)
}
