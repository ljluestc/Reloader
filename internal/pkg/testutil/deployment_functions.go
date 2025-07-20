package testutil

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeleteDeployment handles both string names and Secret objects
func DeleteDeployment(client interface{}, nameOrObj interface{}) error {
	var name string
	var namespace string = "default"

	// Extract name from Secret if given
	if secret, ok := nameOrObj.(*v1.Secret); ok {
		name = secret.Name
		namespace = secret.Namespace
	} else if nameStr, ok := nameOrObj.(string); ok {
		name = nameStr
	}

	// Try to use client if it's a kubernetes client
	if k8sClient, ok := client.(kubernetes.Interface); ok {
		return k8sClient.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	}

	// Dummy implementation for tests
	return nil
}

// DeleteSecret handles both string names and Secret objects
func DeleteSecret(client interface{}, nameOrObj interface{}) error {
	var name string
	var namespace string = "default"

	// Extract name from Secret if given
	if secret, ok := nameOrObj.(*v1.Secret); ok {
		name = secret.Name
		namespace = secret.Namespace
	} else if nameStr, ok := nameOrObj.(string); ok {
		name = nameStr
	}

	// Try to use client if it's a kubernetes client
	if k8sClient, ok := client.(kubernetes.Interface); ok {
		return k8sClient.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	}

	// Dummy implementation for tests
	return nil
}
