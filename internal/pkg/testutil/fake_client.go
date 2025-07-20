package testutil

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// GetFakeKubernetesClient returns a fake Kubernetes client for testing
func GetFakeKubernetesClient() kubernetes.Interface {
	return fake.NewSimpleClientset()
}

// CreateFakeClientSet creates a fake client set for testing
func CreateFakeClientSet() (kubernetes.Interface, string, error) {
	client := fake.NewSimpleClientset()
	namespace := "test-namespace-" + GenerateRandomSequence(5)
	return client, namespace, nil
}
