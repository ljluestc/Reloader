package testutil

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ConvertResourceToSHAString is a helper function to convert a resource to a SHA for testing
func ConvertResourceToSHAString(resourceType, namespace, name string, data interface{}) string {
	// This is a mock implementation that returns a deterministic "SHA" for testing
	return "sha256-" + resourceType + "-" + namespace + "-" + name
}

// DeleteSecret deletes a secret with proper error handling for different argument types
func DeleteSecretByName(client kubernetes.Interface, namespace, name string) error {
	return client.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// DeleteDeploymentByName deletes a deployment with proper error handling
func DeleteDeploymentByName(client kubernetes.Interface, namespace, name string) error {
	return client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
