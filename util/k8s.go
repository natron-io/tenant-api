package util

import (
	"context"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	Clientset      *kubernetes.Clientset
	TENANT_LABEL   string
	DISCOUNT_LABEL string
)

// GetPodsByTenant returns a map of pods for each tenant
func GetPodsByTenant(tenants []string) (map[string][]string, error) {
	tenantPods := make(map[string][]string)
	for _, tenant := range tenants {
		tenantNamespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		// for each tenantNamespace get pods
		for _, namespace := range tenantNamespaces.Items {
			pods, err := Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: TENANT_LABEL + "=" + tenant,
			})
			if err != nil {
				Status = "Error"
				return nil, err
			}

			// for each pod add it to the list of pods for the tenant
			for _, pod := range pods.Items {
				tenantPods[tenant] = append(tenantPods[tenant], pod.Name)
			}
		}
	}
	return tenantPods, nil
}

// GetNamespacesByTenant returns a map of namespaces for each tenant
func GetNamespacesByTenant(tenants []string) (map[string][]string, error) {
	tenantNamespaces := make(map[string][]string)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		// for each tenantNamespace get pods
		for _, namespace := range namespaces.Items {
			tenantNamespaces[tenant] = append(tenantNamespaces[tenant], namespace.Name)
		}
	}
	return tenantNamespaces, nil
}

// GetServiceAccountsByTenant returns a map of service accounts for each tenant
func GetServiceAccountsByTenant(tenants []string) (map[string]map[string][]string, error) {
	tenantServiceAccounts := make(map[string]map[string][]string)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		// for each tenantNamespace get service accounts
		for _, namespace := range namespaces.Items {
			serviceAccounts, err := Clientset.CoreV1().ServiceAccounts(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				Status = "Error"
				return nil, err
			}

			// for each service account add it to the list of service accounts for the namespace
			tenantServiceAccounts[tenant] = make(map[string][]string)
			tenantServiceAccounts[tenant][namespace.Name] = make([]string, 0)
			for _, serviceAccount := range serviceAccounts.Items {
				tenantServiceAccounts[tenant][namespace.Name] = append(tenantServiceAccounts[tenant][namespace.Name], serviceAccount.Name)
			}
		}
	}
	return tenantServiceAccounts, nil
}

// GetCPURequestsSumByTenant returns the sum of CPU requests for each tenant
func GetCPURequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantCPURequests := make(map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		for _, namespace := range namespaces.Items {
			pods, err := Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: TENANT_LABEL + "=" + tenant,
			})

			if err != nil {
				Status = "Error"
				return nil, err
			}

			for _, pod := range pods.Items {

				// get DISCOUNT_REQUEST by DISCOUNT_LABEL
				discount := pod.Labels[DISCOUNT_LABEL]
				if discount == "" {
					discount = "0"
				}
				// convert to float64
				discountFloat, err := strconv.ParseFloat(discount, 64)
				if err != nil || discountFloat < 0 || discountFloat > 1 {
					WarningLogger.Printf("Discount value %s is not valid for pod %s with label %s", discount, pod.Name, DISCOUNT_LABEL)
					discount = "0"
				}

				CPU_DISCOUNT_PERCENT = discountFloat

				tenantCPURequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
			}
		}
	}
	return tenantCPURequests, nil
}

// GetMemoryRequestsSumByTenant returns the sum of memory requests for each tenant
func GetMemoryRequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantMemoryRequests := make(map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		for _, namespace := range namespaces.Items {
			pods, err := Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: TENANT_LABEL + "=" + tenant,
			})
			if err != nil {
				Status = "Error"
				return nil, err
			}

			for _, pod := range pods.Items {

				// get DISCOUNT_REQUEST by DISCOUNT_LABEL
				discount := pod.Labels[DISCOUNT_LABEL]
				if discount == "" {
					discount = "0"
				}
				// convert to float64
				discountFloat, err := strconv.ParseFloat(discount, 64)
				if err != nil || discountFloat < 0 || discountFloat > 1 {
					WarningLogger.Printf("Discount value %s is not valid for pod %s with label %s", discount, pod.Name, DISCOUNT_LABEL)
					discount = "0"
				}

				MEMORY_DISCOUNT_PERCENT = discountFloat

				tenantMemoryRequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Memory().Value()
			}
		}
	}
	return tenantMemoryRequests, nil
}

// GetStorageRequestsSumByTenant returns the sum of storage requests for each tenant
func GetStorageRequestsSumByTenant(tenants []string) (map[string]map[string]int64, error) {
	tenantPVCs := make(map[string]map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		for _, namespace := range namespaces.Items {
			pvcList, err := Clientset.CoreV1().PersistentVolumeClaims(namespace.Name).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				Status = "Error"
				return nil, err
			}

			// create a map for each storage class with a count of pvc size if it exists
			tenantPVCs[tenant] = make(map[string]int64)
			for _, pvc := range pvcList.Items {
				// get DISCOUNT_REQUEST by DISCOUNT_LABEL
				discount := pvc.Labels[DISCOUNT_LABEL]
				if discount == "" {
					discount = "0"
				}
				// convert to float64
				discountFloat, err := strconv.ParseFloat(discount, 64)
				if err != nil || discountFloat < 0 || discountFloat > 1 {
					WarningLogger.Printf("Discount value %s is not valid for pod %s with label %s", discount, pvc.Name, DISCOUNT_LABEL)
					discount = "0"
				}

				STORAGE_DISCOUNT_PERCENT = discountFloat
				tenantPVCs[tenant][*pvc.Spec.StorageClassName] += pvc.Spec.Resources.Requests.Storage().Value()
			}

			// if tenant is emtpy remove it from the map
			if len(tenantPVCs[tenant]) == 0 {
				delete(tenantPVCs, tenant)
			}
		}
	}
	return tenantPVCs, nil
}

// GetIngressRequestsSumByTenant returns the sum of ingress requests for each tenant
func GetIngressRequestsSumByTenant(tenants []string) (map[string][]string, error) {
	tenantsIngress := make(map[string][]string)

	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: TENANT_LABEL + "=" + tenant,
		})

		if err != nil {
			Status = "Error"
			return nil, err
		}

		// get ingress for each namespace in the tenant and add it to the map of ingress for the tenant
		for _, namespace := range namespaces.Items {
			ingressList, err := Clientset.NetworkingV1().Ingresses(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: TENANT_LABEL + "=" + tenant,
			})

			if err != nil {
				Status = "Error"
				return nil, err
			}

			for _, ingress := range ingressList.Items {
				// get DISCOUNT_REQUEST by DISCOUNT_LABEL
				discount := ingress.Labels[DISCOUNT_LABEL]
				if discount == "" {
					discount = "0"
				}
				// convert to float64
				discountFloat, err := strconv.ParseFloat(discount, 64)
				if err != nil || discountFloat < 0 || discountFloat > 1 {
					WarningLogger.Printf("Discount value %s is not valid for pod %s with label %s", discount, ingress.Name, DISCOUNT_LABEL)
					discount = "0"
				}

				INGRESS_DISCOUNT_PERCENT = discountFloat

				// apend ingress hostname to the list of ingress for the tenant
				tenantsIngress[tenant] = append(tenantsIngress[tenant], ingress.Name)
			}
		}
	}

	return tenantsIngress, nil
}
