package testutil

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UpdateSecretData updates an existing secret's data
func UpdateSecretData(client kubernetes.Interface, namespace, name string, data map[string][]byte) (*v1.Secret, error) {
	// Get the existing secret
	secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Update the data
	secret.Data = data

	// Update the secret in the clientset
	updatedSecret, err := client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	return updatedSecret, err
}

// UpdateSecretWithInterface with compatibility for both signatures:
// 1. Controller tests: (*v1.Secret, namespace, name, key, value) -> *v1.Secret
// 2. Other tests: (kubernetes.Interface, namespace, name, key, value) -> (*v1.Secret, error)
func UpdateSecretWithInterface(arg1 interface{}, namespace, name, keyOrData string, valueOrData ...string) (*v1.Secret, error) {
	switch client := arg1.(type) {
	case kubernetes.Interface:
		// New signature: UpdateSecret(client, namespace, name, key, value)
		var data map[string][]byte
		if len(valueOrData) > 0 {
			data = map[string][]byte{
				keyOrData: []byte(valueOrData[0]),
			}
		} else {
			data = map[string][]byte{}
		}

		return UpdateSecretData(client, namespace, name, data)

	case *v1.Secret:
		// Legacy signature: UpdateSecret(secret, namespace, name, key, value) -> *v1.Secret
		secret := client
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		if len(valueOrData) > 0 {
			secret.Data[keyOrData] = []byte(valueOrData[0])
		}
		return secret, nil
	}

	// For controller_test.go, we need to return just the secret
	// This is a dummy implementation that should not be reached in normal flow
	return &v1.Secret{}, nil
}
