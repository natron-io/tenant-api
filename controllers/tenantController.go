package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	LABELSELECTOR string
)

func GetPods(c *fiber.Ctx) error {

	tenants := CheckAuth(c)
	if tenants == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// create a map for each tenant with a list of pods with labels in it
	tenantPods := make(map[string][]string)
	for _, tenant := range tenants {
		tenantNamespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		// for each tenantNamespace get pods
		for _, namespace := range tenantNamespaces.Items {
			pods, err := util.Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: LABELSELECTOR + "=" + tenant,
			})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			// for each pod add it to the list of pods for the tenant
			for _, pod := range pods.Items {
				tenantPods[tenant] = append(tenantPods[tenant], pod.Name)
			}
		}
	}
	util.InfoLogger.Println("/api/v1/pods hit from IP: " + c.IP())
	return c.JSON(tenantPods)
}

func GetNamespaces(c *fiber.Ctx) error {

	tenants := CheckAuth(c)
	if tenants == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// create a map for each tenant with a list of namespaces with labels in it
	tenantNamespaces := make(map[string][]string)
	for _, tenant := range tenants {
		namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		// for each tenantNamespace get pods
		for _, namespace := range namespaces.Items {
			tenantNamespaces[tenant] = append(tenantNamespaces[tenant], namespace.Name)
		}
	}
	util.InfoLogger.Println("/api/v1/namespaces hit from IP: " + c.IP())
	return c.JSON(tenantNamespaces)
}

func GetServiceAccounts(c *fiber.Ctx) error {
	tenants := CheckAuth(c)
	if tenants == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	// create a map for each tenant with a map of namespaces with a list of service accounts with labels in it
	tenantServiceAccounts := make(map[string]map[string][]string)
	for _, tenant := range tenants {
		namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		// for each tenantNamespace get service accounts
		for _, namespace := range namespaces.Items {
			serviceAccounts, err := util.Clientset.CoreV1().ServiceAccounts(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			// for each service account add it to the list of service accounts for the namespace
			tenantServiceAccounts[tenant] = make(map[string][]string)
			tenantServiceAccounts[tenant][namespace.Name] = make([]string, 0)
			for _, serviceAccount := range serviceAccounts.Items {
				tenantServiceAccounts[tenant][namespace.Name] = append(tenantServiceAccounts[tenant][namespace.Name], serviceAccount.Name)
			}
		}
	}
	util.InfoLogger.Println("/api/v1/serviceaccounts hit from IP: " + c.IP())
	return c.JSON(tenantServiceAccounts)
}

func GetCPURequestsSum(c *fiber.Ctx) error {
	tenants := CheckAuth(c)
	if tenants == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// create a map for each tenant with a added cpu requests
	tenantCPURequests := make(map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		for _, namespace := range namespaces.Items {
			pods, err := util.Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: LABELSELECTOR + "=" + tenant,
			})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			for _, pod := range pods.Items {
				tenantCPURequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
			}
		}
	}

	util.InfoLogger.Println("/api/v1/cpurequests hit from IP: " + c.IP())
	return c.JSON(tenantCPURequests)
}

func GetMemoryRequestsSum(c *fiber.Ctx) error {
	tenants := CheckAuth(c)
	if tenants == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// create a map for each tenant with a added memory requests
	tenantMemoryRequests := make(map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		for _, namespace := range namespaces.Items {
			pods, err := util.Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: LABELSELECTOR + "=" + tenant,
			})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			for _, pod := range pods.Items {
				// byte to megabyte
				tenantMemoryRequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024
			}
		}
	}

	util.InfoLogger.Println("/api/v1/memoryrequests hit from IP: " + c.IP())
	return c.JSON(tenantMemoryRequests)
}

func GetStorageAllocationSum(c *fiber.Ctx) error {
	tenants := CheckAuth(c)
	if tenants == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// create a map for each tenant with a map of storage classes with calculated pvcs in it
	tenantPVCs := make(map[string]map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := util.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		for _, namespace := range namespaces.Items {
			pvcList, err := util.Clientset.CoreV1().PersistentVolumeClaims(namespace.Name).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			// create a map for each storage class with a count of pvc size
			tenantPVCs[tenant] = make(map[string]int64)
			for _, pvc := range pvcList.Items {
				tenantPVCs[tenant][*pvc.Spec.StorageClassName] += pvc.Spec.Resources.Requests.Storage().Value()
			}
		}
	}

	util.InfoLogger.Println("/api/v1/pvcs hit from IP: " + c.IP())
	return c.JSON(tenantPVCs)

}
