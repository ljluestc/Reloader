package testutil_test

import (
	"os"
	"testing"

	"github.com/stakater/Reloader/internal/pkg/testutil"
)

// No TestMain needed here since the main test setup is now in test_environment.go

// Setup additional environment variables for specific tests
func setupAdditionalTestEnv() {
	os.Setenv("POD_NAME", "test-pod")
	os.Setenv("POD_NAMESPACE", "test-namespace")
}

// Cleanup additional environment variables for specific tests
func cleanupAdditionalTestEnv() {
	os.Unsetenv("POD_NAME")
	os.Unsetenv("POD_NAMESPACE")
}

// Example test that needs additional environment setup
func TestWithAdditionalEnv(t *testing.T) {
	setupAdditionalTestEnv()
	defer cleanupAdditionalTestEnv()

	// Test logic here
}
