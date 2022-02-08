package util

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetRessourceQuota returns the resource quota for the given tenant and label set in the config namespace
func GetRessourceQuota(tenant string) (v1.ResourceQuota, error) {
	// get resource quota in namespace "test-tenant-config"

	// get resource quota from namespace
	quota, err := Clientset.CoreV1().ResourceQuotas(tenant).Get(context.TODO(), tenant, metav1.GetOptions{})
	if err != nil {
		return v1.ResourceQuota{}, err
	}

	return *quota, nil
}
