package helper

func IsServiceAllowed(serviceName string) bool {
	allowedServices := map[string]bool{
		"order_service": true,
	}

	if allowedServices[serviceName] {
		return true
	}

	return false
}
