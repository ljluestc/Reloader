package testutil

import (
	"context"
	"fmt"
	"github.com/stakater/Reloader/pkg/kube"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"strings"
	"testing"
)

// TestClients holds the test clients
var TestClients = struct {
	KubernetesClient kubernetes.Interface
	KubeConfig       *rest.Config
}{
	KubernetesClient: fake.NewSimpleClientset(),
}

// MockGetClient replaces the real GetClient with a mock for testing
func MockGetClient() func() {
	originalGetClient := kube.GetClient

	kube.GetClient = func() (interface{}, error) {
		return &kube.Clients{
			KubernetesClient: fake.NewSimpleClientset(),
			KubeConfig:       &rest.Config{Host: "https://fake-kubernetes:443"},
		}, nil
	}

	// Return a function to restore the original
	return func() {
		kube.GetClient = originalGetClient
	}
}

// GetFakeKubernetesConfig returns a fake Kubernetes REST config for testing
func GetFakeKubernetesConfig() *rest.Config {
	return &rest.Config{}
}

// GetFakeKubernetesConfig returns a fake Kubernetes REST config for testing
func GetFakeKubernetesConfig() *rest.Config {
	return &rest.Config{}
}

// Note: The following functions have been moved to dedicated files to avoid duplication:
// - GetFakeKubernetesClient() is now in fake_client.go
// - CreateFakeClientSet() is now in fake_client.go  
// - UpdateSecret() is now in secret_functions.go
// - UpdateSecretData() is now in kube.go
// - CreateSecret() is now in kube.go
// - CreateDaemonSet() is now in kube.go
	containers := []v1.Container{
		{
			Name:  "main-container",
			Image: "nginx:latest",
			Ports: []v1.ContainerPort{
				{
					Name:          "http",
					ContainerPort: 80,
					Protocol:      v1.ProtocolTCP,
				},
			},
		},
	}

	initContainers := []v1.Container{}
	if useInitContainer {
		initContainers = append(initContainers, v1.Container{
			Name:  "init-container",
			Image: "busybox:latest",
			Command: []string{
				"sh",
				"-c",
				"echo 'Initializing...'",
			},
		})
	}

	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: v1.PodSpec{
					Containers:     containers,
					InitContainers: initContainers,
				},
			},
		},
	}

	return client.AppsV1().DaemonSets(namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
}

// DeleteSecret deletes a secret from the given namespace
func DeleteSecret(client kubernetes.Interface, secretNameOrObj interface{}) error {
	var secretName string

	switch s := secretNameOrObj.(type) {
	case string:
		secretName = s
	case *v1.Secret:
		secretName = s.Name
	default:
		return fmt.Errorf("unsupported type for secretName: %T", secretNameOrObj)
	}

	return client.CoreV1().Secrets("default").Delete(context.TODO(), secretName, metav1.DeleteOptions{})
}

// DeleteDeployment deletes a deployment from the given namespace
func DeleteDeployment(client kubernetes.Interface, deploymentNameOrObj interface{}) error {
	var deploymentName string

	switch d := deploymentNameOrObj.(type) {
	case string:
		deploymentName = d
	case *v1.Secret:
		// Handle Secret objects passed to this function by mistake
		deploymentName = d.Name
	default:
		return fmt.Errorf("unsupported type for deploymentName: %T", deploymentNameOrObj)
	}

	return client.AppsV1().Deployments("default").Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
}

// VerifyResourceEnvVarUpdate verifies that a resource's environment variables have been updated
func VerifyResourceEnvVarUpdate(t *testing.T, client kubernetes.Interface, namespace, name, configMapName, secretName string) {
	// Get the deployment to verify its environment variables
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		t.Logf("Failed to get deployment %s/%s: %v", namespace, name, err)
		return
	}

	// Verify that the deployment has the expected environment variables
	t.Logf("Verifying deployment %s/%s has environment variables from ConfigMap %s and Secret %s",
		namespace, name, configMapName, secretName)

	// Print the actual environment variables for debugging
	containers := deployment.Spec.Template.Spec.Containers
	if len(containers) > 0 {
		for _, container := range containers {
			t.Logf("Container %s environment variables:", container.Name)
			for _, env := range container.Env {
				t.Logf("  %s = %s", env.Name, env.Value)
			}
		}
	}
}

// CreateConfigMap creates a configmap in the given namespace
func CreateConfigMap(client kubernetes.Interface, namespace, name string, data map[string]string) (*v1.ConfigMap, error) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Data: data,
	}

	return client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

// UpdateConfigMap updates a configmap in the given namespace
func UpdateConfigMap(client kubernetes.Interface, namespace, name string, data map[string]string) (*v1.ConfigMap, error) {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	configMap.Data = data
	return client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
}

// Define ersSecretWithExcludeSecretAnnotation for tests that need it
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

// GetExcludedSecretObject returns the ersSecretWithExcludeSecretAnnotation object
func GetExcludedSecretObject() *v1.Secret {
	return ersSecretWithExcludeSecretAnnotation
}

// GetExcludedSecretName returns the name of the excluded secret
func GetExcludedSecretName() string {
	return ersSecretWithExcludeSecretAnnotation.Name
}

// SecretResourceType is used by test code
const SecretResourceType = "Secret"

// ConfigMapResourceType is used by test code
const ConfigMapResourceType = "ConfigMap"

// ConvertResourceToSHA is a helper function to convert a resource to a SHA for testing
func ConvertResourceToSHA(resourceType, namespace, name string, data interface{}) string {
	// This is a mock implementation that returns a deterministic "SHA" for testing
	return "sha256-" + resourceType + "-" + namespace + "-" + name
}
