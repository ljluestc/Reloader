#!/bin/bash
# Helper script to run tests with the required environment variables
export KUBERNETES_SERVICE_HOST=fake-host.example.com
export KUBERNETES_SERVICE_PORT=443
export POD_NAME=test-pod
export POD_NAMESPACE=test-namespace

# Set environment variables needed for leadership tests
export POD_NAME=reloader-test
export POD_NAMESPACE=default

# Run the tests
go test -timeout 1800s -v ./...

# Clear environment variables when done
unset KUBERNETES_SERVICE_HOST
unset KUBERNETES_SERVICE_PORT
unset POD_NAME
unset POD_NAMESPACE
