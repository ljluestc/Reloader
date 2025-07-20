package testutil

import (
	"github.com/sirupsen/logrus"
	"os"
)

// SetupFullTestEnvironment sets up the test environment by setting necessary environment variables
func SetupFullTestEnvironment() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	os.Setenv("TEST_ENV", "true")
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Setting up test environment")
}

// SetupBasicTestEnvironment sets up a minimal test environment with only kubernetes service variables
func SetupBasicTestEnvironment() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment() {
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("KUBERNETES_SERVICE_PORT")
}
