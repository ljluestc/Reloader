package handler

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This file exists to provide necessary imports for the test files
// It ensures v1.ConfigMap and metav1 are available to upgrade_test.go
