package util

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	Clientset     *kubernetes.Clientset
	LABELSELECTOR string
)

func GetPodsByTenant(tenants []string) (map[string][]string, error) {
	tenantPods := make(map[string][]string)
	for _, tenant := range tenants {
		tenantNamespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return nil, err
		}

		// for each tenantNamespace get pods
		for _, namespace := range tenantNamespaces.Items {
			pods, err := Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: LABELSELECTOR + "=" + tenant,
			})
			if err != nil {
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

func GetNamespacesByTenant(tenants []string) (map[string][]string, error) {
	tenantNamespaces := make(map[string][]string)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return nil, err
		}

		// for each tenantNamespace get pods
		for _, namespace := range namespaces.Items {
			tenantNamespaces[tenant] = append(tenantNamespaces[tenant], namespace.Name)
		}
	}
	return tenantNamespaces, nil
}

func GetServiceAccountsByTenant(tenants []string) (map[string]map[string][]string, error) {
	tenantServiceAccounts := make(map[string]map[string][]string)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return nil, err
		}

		// for each tenantNamespace get service accounts
		for _, namespace := range namespaces.Items {
			serviceAccounts, err := Clientset.CoreV1().ServiceAccounts(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
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

func GetCPURequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantCPURequests := make(map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return nil, err
		}

		for _, namespace := range namespaces.Items {
			pods, err := Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: LABELSELECTOR + "=" + tenant,
			})
			if err != nil {
				return nil, err
			}

			for _, pod := range pods.Items {
				tenantCPURequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
			}
		}
	}
	return tenantCPURequests, nil
}

func GetMemoryRequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantMemoryRequests := make(map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return nil, err
		}

		for _, namespace := range namespaces.Items {
			pods, err := Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: LABELSELECTOR + "=" + tenant,
			})
			if err != nil {
				return nil, err
			}

			for _, pod := range pods.Items {
				tenantMemoryRequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Memory().Value()
			}
		}
	}
	return tenantMemoryRequests, nil
}

func GetStorageRequestsSumByTenant(tenants []string) (map[string]map[string]int64, error) {
	tenantPVCs := make(map[string]map[string]int64)
	for _, tenant := range tenants {
		namespaces, err := Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
			LabelSelector: LABELSELECTOR + "=" + tenant,
		})

		if err != nil {
			return nil, err
		}

		for _, namespace := range namespaces.Items {
			pvcList, err := Clientset.CoreV1().PersistentVolumeClaims(namespace.Name).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return nil, err
			}

			// create a map for each storage class with a count of pvc size
			tenantPVCs[tenant] = make(map[string]int64)
			for _, pvc := range pvcList.Items {
				tenantPVCs[tenant][*pvc.Spec.StorageClassName] += pvc.Spec.Resources.Requests.Storage().Value()
			}
		}
	}
	return tenantPVCs, nil
}
