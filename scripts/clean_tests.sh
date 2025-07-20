#!/bin/bash

# Script to clean and fix test files issues
# This script will detect problematic files and fix common issues

echo "Checking for test files with issues..."

# Find all test files with package declaration issues
for file in $(find . -name "*_test.go"); do
  # Check if the file has duplicate package declarations
  PACKAGE_COUNT=$(grep -c "^package " $file)
  
  if [ $PACKAGE_COUNT -gt 1 ]; then
    echo "Fixing duplicate package declarations in $file"
    # Keep only the first package declaration
    sed -i '0,/^package /!s/^package .*//' $file
  fi
  
  # Check if the file has test package but doesn't end with _test
  if grep -q "package.*_test" $file && [[ ! $(basename $file) =~ "_test.go" ]]; then
    echo "Warning: $file has a test package but doesn't follow naming convention"
  fi
done

# Find all files with duplicate import blocks
for file in $(find . -name "*.go"); do
  # Count import blocks
  IMPORT_COUNT=$(grep -c "^import (" $file)
  
  if [ $IMPORT_COUNT -gt 1 ]; then
    echo "Warning: $file has multiple import blocks, consider merging them"
  fi
done

# Fix test environment setup
for testdir in $(find . -type d -name "testutil"); do
  if [ ! -f "$testdir/env.go" ]; then
    echo "Creating environment setup in $testdir"
    cat > "$testdir/env.go" << EOF
package testutil

import (
	"os"
)

// SetupTestEnv sets up the test environment by setting necessary environment variables
func SetupTestEnv() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "fake-host.example.com")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
}

// TeardownTestEnv cleans up the test environment
func TeardownTestEnv() {
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("KUBERNETES_SERVICE_PORT")
}
EOF
  fi
done

echo "Cleanup complete!"
