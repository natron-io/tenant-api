package util

// Contains returns true if the provided string is in the provided slice
func Contains(tenant string, tenants []string) bool {
	for _, t := range tenants {
		if t == tenant {
			return true
		}
	}
	return false
}
