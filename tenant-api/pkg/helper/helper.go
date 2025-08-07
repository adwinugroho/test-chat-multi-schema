package helper

import "strings"

func SanitizeTenantID(tenantID string) string {
	return strings.ReplaceAll(tenantID, "-", "_")
}
