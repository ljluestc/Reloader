package testutil

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Constants for resource types
const (
	SecretResourceType    = "Secret"
	ConfigMapResourceType = "ConfigMap"
)

// Create a standard excluded secret for tests
var ersSecretWithExcludeSecretAnnotation = &v1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "excluded-secret",
		Namespace: "default",
		Annotations: map[string]string{
			"reloader.stakater.com/exclude-secret": "true",
		},
	},
	Data: map[string][]byte{
		"key1": []byte("value1"),
	},
	Type: v1.SecretTypeOpaque,
}

// DeleteSecretDeployment for handler tests that pass a Secret object
func DeleteSecretDeployment(client interface{}, nameOrObj interface{}) error {
	return nil // Mock implementation for tests
}
