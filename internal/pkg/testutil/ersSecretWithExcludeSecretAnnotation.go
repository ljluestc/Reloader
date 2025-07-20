package testutil

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Define the excluded secret for tests
var ExcludedSecret = &v1.Secret{
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
