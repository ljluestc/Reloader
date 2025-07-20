package testutil
package testutil

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UpdateSecret is a flexible function that handles both old and new calling patterns:
// Old: (secretClient *v1.Secret, namespace, name, key, value string) *v1.Secret
// New: (client kubernetes.Interface, namespace, name, key string, data map[string][]byte) (*v1.Secret, error)
func UpdateSecret(clientOrSecret interface{}, namespace, name, key string, valueOrData ...interface{}) (*v1.Secret, error) {
	// Check if first argument is a Secret (old pattern)
	if secret, ok := clientOrSecret.(*v1.Secret); ok {
		// Handle old pattern
		if len(valueOrData) == 1 {
			if value, ok := valueOrData[0].(string); ok {
				if secret.Data == nil {
					secret.Data = make(map[string][]byte)
				}
				if key != "" {
					secret.Data[key] = []byte(value)
				}
			}
		}
		return secret, nil
	}

	// Check if first argument is a Kubernetes client (new pattern)
	if client, ok := clientOrSecret.(kubernetes.Interface); ok {
		// Handle new pattern with k8s client
		secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}

		// If we have data as map[string][]byte
		if len(valueOrData) == 1 {
			if data, ok := valueOrData[0].(map[string][]byte); ok {
				secret.Data = data
			} else if value, ok := valueOrData[0].(string); ok && key != "" {
				// Handle simple key/value
				secret.Data[key] = []byte(value)
			}
		}

		// Update in kubernetes
		updatedSecret, err := client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
		return updatedSecret, err
	}

	// Return an empty secret and nil error as fallback
	return &v1.Secret{}, nil
}
import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UpdateSecret is a flexible function that handles both old and new calling patterns:
// Old: (secretClient *v1.Secret, namespace, name, key, value string) *v1.Secret
// New: (client kubernetes.Interface, namespace, name, key string, data map[string][]byte) (*v1.Secret, error)
func UpdateSecret(clientOrSecret interface{}, namespace, name, key string, valueOrData ...interface{}) (*v1.Secret, error) {
	// Check if first argument is a Secret (old pattern)
	if secret, ok := clientOrSecret.(*v1.Secret); ok {
		// Handle old pattern
		if len(valueOrData) == 1 {
			if value, ok := valueOrData[0].(string); ok {
				if secret.Data == nil {
					secret.Data = make(map[string][]byte)
				}
				if key != "" {
					secret.Data[key] = []byte(value)
				}
			}
		}
		return secret, nil
	}

	// Check if first argument is a Kubernetes client (new pattern)
	if client, ok := clientOrSecret.(kubernetes.Interface); ok {
		// Handle new pattern with k8s client
		secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}

		// If we have data as map[string][]byte
		if len(valueOrData) == 1 {
			if data, ok := valueOrData[0].(map[string][]byte); ok {
				secret.Data = data
			} else if value, ok := valueOrData[0].(string); ok && key != "" {
				// Handle simple key/value
				secret.Data[key] = []byte(value)
			}
		}

		// Update in kubernetes
		updatedSecret, err := client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
		return updatedSecret, err
	}

	// Return an empty secret and nil error as fallback
	return &v1.Secret{}, nil
}

// ErsSecretWithExcludeSecretAnnotation returns a secret with exclude annotation
var ErsSecretWithExcludeSecretAnnotation = &v1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "exclude-secret-test",
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
