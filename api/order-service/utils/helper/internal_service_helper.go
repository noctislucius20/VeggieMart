package helper

func IsServiceAllowed(serviceName string) bool {
	allowedServices := map[string]bool{
		"payment_service": true,
	}

	if allowedServices[serviceName] {
		return true
	}

	return false
}
