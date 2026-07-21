package helper

import "strings"

func HasRequiredPermissions(rolePermissions []string, requiredPermissions []string) bool {
	if len(rolePermissions) == 0 {
		return false
	}

	permissionSet := map[string]struct{}{}

	for _, perm := range rolePermissions {
		perm = strings.TrimSpace(perm)
		permissionSet[perm] = struct{}{}
	}

	for _, required := range requiredPermissions {
		if _, exists := permissionSet[required]; !exists {
			return false
		}
	}

	return true
}
