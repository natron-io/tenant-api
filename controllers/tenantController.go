package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(c *fiber.Ctx) error {
	// Get all pods in the cluster
	pods, err := util.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	// // Get all pods in all namespaces by label
	// pods, err := util.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: "app=tenant-api"})

	// // Get pvc by label
	// pvc, err := util.Clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{LabelSelector: "tenant=tenant-api"})

	// // Get cpu request by label
	// cpuRequest, err := util.Clientset.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{LabelSelector: "tenant=tenant-api"})

	util.InfoLogger.Println("/api/v1/pods hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}

	return c.JSON(pods.Items)

}

func GetNamespaces(c *fiber.Ctx) error {
	namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	util.InfoLogger.Println("/api/v1/namespaces hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}

	return c.JSON(namespaces.Items)
}

func GetServiceAccounts(c *fiber.Ctx) error {
	serviceAccounts, err := util.Clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})

	util.InfoLogger.Println("/api/v1/serviceAccounts hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}

	return c.JSON(serviceAccounts.Items)
}
