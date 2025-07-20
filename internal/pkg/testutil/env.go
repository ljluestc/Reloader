package testutil
package testutil

import (
	"os"
	"github.com/sirupsen/logrus"
)

// SetupTestEnvironment sets up the test environment by setting necessary environment variables
func SetupTestEnvironment() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	os.Setenv("TEST_ENV", "true")
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Setting up test environment")
}
import (
	"os"
)

// SetupTestEnvironment sets up the test environment by setting necessary environment variables
func SetupTestEnvironment() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment() {
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("KUBERNETES_SERVICE_PORT")
}
