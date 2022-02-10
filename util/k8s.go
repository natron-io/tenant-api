package util

import (
	"context"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	Clientset       *kubernetes.Clientset
	DISCOUNT_LABEL  string
	EXCLUDE_STRINGS []string
)

// GetPodsByTenant returns a map of pods for each tenant
func GetPodsByTenant(tenants []string) (map[string][]string, error) {
	tenantPods := make(map[string][]string)
	// get namespace with same name as tenant and get pods
	for _, tenant := range tenants {
		pods, err := Clientset.CoreV1().Pods(tenant).List(context.TODO(), metav1.ListOptions{})
		if err != nil && !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// for each pod add it to the list of pods for the namespace
		tenantPods[tenant] = make([]string, 0)
		for _, pod := range pods.Items {
			tenantPods[tenant] = append(tenantPods[tenant], strings.Split(pod.Name, "-x-")[0])
		}
	}

	return tenantPods, nil
}

func GetPVCsByTenantByStorageClass(tenants []string) (map[string]map[string][]string, error) {
	tenantPVCs := make(map[string]map[string][]string)
	storageClasses, err := GetStorageClassesInCluster()
	if err != nil {
		return nil, err
	}

	for _, tenant := range tenants {
		tenantPVCs[tenant] = make(map[string][]string)
		for _, storageClass := range storageClasses {
			// get pvc by storageclass pvc spec
			pvcs, err := Clientset.CoreV1().PersistentVolumeClaims(tenant).List(context.TODO(), metav1.ListOptions{})
			if err != nil && !strings.Contains(err.Error(), "not found") {
				return nil, err
			}

			for _, pvc := range pvcs.Items {
				if *pvc.Spec.StorageClassName == storageClass {
					tenantPVCs[tenant][storageClass] = append(tenantPVCs[tenant][storageClass], pvc.Name)
				}
			}
		}
	}

	return tenantPVCs, nil
}

// GetCPURequestsSumByTenant returns the sum of CPU requests for each tenant
func GetCPURequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantCPURequests := make(map[string]int64)
	for _, tenant := range tenants {
		pods, err := Clientset.CoreV1().Pods(tenant).List(context.TODO(), metav1.ListOptions{})

		if err != nil && !strings.Contains(err.Error(), "not found") {
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
				return nil, err
			}

			CPU_DISCOUNT_PERCENT = discountFloat

			tenantCPURequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
		}
	}
	return tenantCPURequests, nil
}

// GetMemoryRequestsSumByTenant returns the sum of memory requests for each tenant
func GetMemoryRequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantMemoryRequests := make(map[string]int64)
	for _, tenant := range tenants {
		pods, err := Clientset.CoreV1().Pods(tenant).List(context.TODO(), metav1.ListOptions{})
		if err != nil && !strings.Contains(err.Error(), "not found") {
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
				return nil, err
			}

			MEMORY_DISCOUNT_PERCENT = discountFloat

			tenantMemoryRequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Memory().Value()
		}
	}
	return tenantMemoryRequests, nil
}

// GetStorageRequestsSumByTenant returns the sum of storage requests for each tenant
func GetStorageRequestsSumByTenant(tenants []string) (map[string]map[string]int64, error) {
	tenantPVCs := make(map[string]map[string]int64)
	for _, tenant := range tenants {
		pvcList, err := Clientset.CoreV1().PersistentVolumeClaims(tenant).List(context.TODO(), metav1.ListOptions{})

		if err != nil && !strings.Contains(err.Error(), "not found") {
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
				return nil, err
			}

			STORAGE_DISCOUNT_PERCENT = discountFloat
			tenantPVCs[tenant][*pvc.Spec.StorageClassName] += pvc.Spec.Resources.Requests.Storage().Value()
		}

		// if tenant is emtpy remove it from the map
		if len(tenantPVCs[tenant]) == 0 {
			delete(tenantPVCs, tenant)
		}
	}
	return tenantPVCs, nil
}

// GetIngressRequestsSumByTenant returns the sum of ingress requests for each tenant
func GetIngressRequestsSumByTenant(tenants []string) (map[string][]string, error) {
	tenantsIngress := make(map[string][]string)

	for _, tenant := range tenants {
		// get ingress for each namespace in the tenant and add it to the map of ingress for the tenant
		ingressList, err := Clientset.NetworkingV1().Ingresses(tenant).List(context.TODO(), metav1.ListOptions{})

		if err != nil && !strings.Contains(err.Error(), "not found") {
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
				return nil, err
			}

			INGRESS_DISCOUNT_PERCENT = discountFloat

			if strings.Contains(ingress.Name, "vcluster") && EXCLUDE_INGRESS_VCLUSTER {
				continue
			}

			// apend ingress hostname to the list of ingress for the tenant
			for _, rule := range ingress.Spec.Rules {
				tenantsIngress[tenant] = append(tenantsIngress[tenant], rule.Host)
			}
		}
	}

	return tenantsIngress, nil
}

func GetStorageClassesInCluster() ([]string, error) {
	storageClasses := make([]string, 0)
	scList, err := Clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})

	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	for _, sc := range scList.Items {
		storageClasses = append(storageClasses, sc.Name)
	}

	return storageClasses, nil
}
