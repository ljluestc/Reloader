package testutil

import (
	"testing"
)

// Simple test to ensure the package can be loaded
func TestPackageLoad(t *testing.T) {
	client := GetFakeKubernetesClient()
	if client == nil {
		t.Error("Failed to create fake client")
	}
}
