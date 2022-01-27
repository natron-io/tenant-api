package util

func Contains(tenant string, tenants []string) bool {
	for _, t := range tenants {
		if t == tenant {
			return true
		}
	}
	return false
}
