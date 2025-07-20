package testutil

import (
	"os"
	"testing"
)

// TestMain is the common entry point for all tests in this package
func TestMain(m *testing.M) {
	// Setup test environment
	SetupTestEnv()

	// Run the tests
	exitCode := m.Run()

	// Cleanup test environment
	CleanupTestEnv()

	// Exit with the test exit code
	os.Exit(exitCode)
}
