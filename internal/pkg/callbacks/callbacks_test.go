package callbacks_test

import (
	"github.com/stakater/Reloader/internal/pkg/testutil"
	"testing"

	"github.com/stakater/Reloader/internal/pkg/crypto"
	"github.com/stakater/Reloader/pkg/kube"
)

func TestSetupCallbacks(t *testing.T) {
	// Setup for this specific test
	testutil.SetupTestEnv()

	// Save original GetClient function
	originalGetClient := kube.GetClient

	// Override with mock function for this test
	kube.GetClient = func() (interface{}, error) {
		client, _, _ := testutil.CreateFakeClientSet()
		return &kube.Clients{
			KubernetesClient: client,
			KubeConfig:       nil,
		}, nil
	}

	// Defer cleanup
	defer func() {
		// Restore original function
		kube.GetClient = originalGetClient

		// Cleanup test environment
		testutil.TeardownTestEnv()
	}()

	// Simple test to ensure the package can be loaded without errors
	sha := crypto.GenerateSHA("test")
	if sha == "" {
		t.Errorf("Failed to generate SHA")
	}
}
