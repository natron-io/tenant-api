package util

// Contains returns true if the provided string is in the provided slice
func Contains(tenant string, tenants []string) bool {
	for _, t := range tenants {
		if t == tenant {
			return true
		}
	}
	if DEBUG {
		InfoLogger.Printf("tenant %s is not in the list of tenants %v", tenant, tenants)
		return true
	}
	return false
}
