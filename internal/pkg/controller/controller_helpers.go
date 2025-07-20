package controller

import (
	v1 "k8s.io/api/core/v1"
)

// UpdateSecretHelper is a wrapper function for tests that expect a single return value
// This maintains compatibility with the old signature used in controller_test.go
func UpdateSecretHelper(secret *v1.Secret, namespace, name, key, value string) *v1.Secret {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[key] = []byte(value)
	return secret
}
