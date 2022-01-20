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

	// Get names of pods
	podNames := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		podNames[i] = pod.Name
	}

	// Return pod names as JSON
	return c.JSON(podNames)

}

func GetPodsByLabel(c *fiber.Ctx) error {
	// Get pods in all the namespaces by label provided
	pods, err := util.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: c.Params("label")})

	util.InfoLogger.Println("/api/v1/pods/label hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}

	// Only return pod name with label
	podNames := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		podNames[i] = pod.Name
	}

	// Return pod names as JSON
	return c.JSON(podNames)
}

func GetNamespaces(c *fiber.Ctx) error {
	namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	util.InfoLogger.Println("/api/v1/namespaces hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}

	// Get names of namespaces
	namespaceNames := make([]string, len(namespaces.Items))
	for i, namespace := range namespaces.Items {
		namespaceNames[i] = namespace.Name
	}

	// Return namespace names as JSON
	return c.JSON(namespaceNames)
}

func GetServiceAccounts(c *fiber.Ctx) error {
	serviceAccounts, err := util.Clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})

	util.InfoLogger.Println("/api/v1/serviceAccounts hit from IP: " + c.IP())
	if err != nil {
		util.WarningLogger.Println(err.Error())
	}

	return c.JSON(serviceAccounts.Items)
}
