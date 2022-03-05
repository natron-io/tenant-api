package util

import (
	"testing"

	"github.com/natron-io/tenant-api/util"
)

func TestContains(t *testing.T) {
	tenants := []string{"tenant1", "tenant2", "tenant3"}
	if !util.Contains("tenant1", tenants) {
		t.Error("tenant1 should be in the tenants slice")
	}
	if util.Contains("tenant4", tenants) {
		t.Error("tenant4 should not be in the tenants slice")
	}
}
