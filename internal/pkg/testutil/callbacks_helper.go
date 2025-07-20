package testutil

import (
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

// This file contains helper functions for callbacks tests
// to avoid import cycles between callbacks and testutil packages

// GetCallbacksTestClient returns a fake client for callbacks tests
func GetCallbacksTestClient() (*fake.Clientset, *rest.Config) {
	client := fake.NewSimpleClientset()
	config := &rest.Config{}
	return client, config
}

// SetupCallbacksTest prepares the environment for callbacks tests
func SetupCallbacksTest() {
	SetupTestEnv()
}

// TeardownCallbacksTest cleans up after callbacks tests
func TeardownCallbacksTest() {
	TeardownTestEnv()
}

// MockCallbacks defines the structure of callback functions used in tests
// without directly importing the callbacks package to avoid import cycles
type MockCallbacks struct {
	ConfigMap   func(namespace, name string)
	Secret      func(namespace, name string)
	Deployment  func(namespace, name string)
	StatefulSet func(namespace, name string)
	DaemonSet   func(namespace, name string)
}

// GetMockCallbacks returns a struct of empty callback functions for testing
func GetMockCallbacks() MockCallbacks {
	return MockCallbacks{
		ConfigMap:   func(namespace, name string) {},
		Secret:      func(namespace, name string) {},
		Deployment:  func(namespace, name string) {},
		StatefulSet: func(namespace, name string) {},
		DaemonSet:   func(namespace, name string) {},
	}
}
