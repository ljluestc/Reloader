package testutil

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestMain is the common entry point for all tests in this package
func TestMain(m *testing.M) {
	// Set up test environment
	SetupTestEnv()

	// Run tests
	code := m.Run()

	// Tear down test environment
	CleanupTestEnv()

	os.Exit(code)
}

// SetupTestEnv sets up the test environment by setting necessary environment variables
func SetupTestEnv() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	os.Setenv("TEST_ENV", "true")
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Setting up test environment")
}

// CleanupTestEnv cleans up the test environment
func CleanupTestEnv() {
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("KUBERNETES_SERVICE_PORT")
	os.Unsetenv("TEST_ENV")
	logrus.Info("Tearing down test environment")
}

// The following are aliases for backward compatibility

// SetupTestEnvironment is an alias for SetupTestEnv
func SetupTestEnvironment() {
	SetupTestEnv()
}

// TeardownTestEnvironment is an alias for CleanupTestEnv
func TeardownTestEnvironment() {
	CleanupTestEnv()
}

// InitTestEnv is an alias for SetupTestEnv
func InitTestEnv() {
	SetupTestEnv()
}

// TeardownTestEnv is an alias for CleanupTestEnv - alias function
func TeardownTestEnv() {
	CleanupTestEnv()
}

// TearDownTestEnv is an alias for CleanupTestEnv - alias with different capitalization
func TearDownTestEnv() {
	CleanupTestEnv()
}

// SetupK8sTestEnv sets up the Kubernetes environment for testing
func SetupK8sTestEnv() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
}

// TeardownK8sTestEnv cleans up the Kubernetes test environment
func TeardownK8sTestEnv() {
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("KUBERNETES_SERVICE_PORT")
}

// CleanupKubernetesTestEnvironment is an alias for TeardownK8sTestEnv
func CleanupKubernetesTestEnvironment() {
	TeardownK8sTestEnv()
}

// Note: GenerateRandomSequence function has been moved to init.go to avoid duplication.
