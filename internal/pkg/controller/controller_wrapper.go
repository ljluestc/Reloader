package controller

import (
	"github.com/stakater/Reloader/internal/pkg/testutil"
	v1 "k8s.io/api/core/v1"
)

// This is a wrapper function to maintain compatibility with the controller tests
// It uses the old signature (secret, namespace, name, key, value) -> secret
func updateSecretWrapper(secret *v1.Secret, namespace, name, key, value string) *v1.Secret {
	return testutil.UpdateSecretLegacy(secret, namespace, name, key, value)
}
