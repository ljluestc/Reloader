package testutil

import (
	v1 "k8s.io/api/core/v1"
)

// GetSecretName returns the name of a secret
// Handles both string names and Secret objects
func GetSecretName(secretNameOrObj interface{}) string {
	switch s := secretNameOrObj.(type) {
	case string:
		return s
	case *v1.Secret:
		return s.Name
	default:
		return ""
	}
}

// UpdateSecretLegacy supports the old signature used in controller_test.go
// This function intentionally returns a single value (Secret) for compatibility
func UpdateSecretLegacy(secret *v1.Secret, namespace, name, key, value string) *v1.Secret {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[key] = []byte(value)
	return secret
}
