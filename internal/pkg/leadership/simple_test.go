package leadership

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// This simple test is in the leadership package itself
// for direct access to unexported functions/variables.
// The _test package tests should focus on integration testing.

func TestHealthzFunction(t *testing.T) {
	// Save the current state of the healthy variable
	originalHealthy := healthy
	
	// Test healthy state
	healthy = true
	
	request, err := http.NewRequest(http.MethodGet, "/live", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	response := httptest.NewRecorder()
	healthz(response, request)
	
	if response.Code != 200 {
		t.Errorf("Expected status code 200, got %d", response.Code)
	}
	
	// Test unhealthy state
	healthy = false
	
	request, err = http.NewRequest(http.MethodGet, "/live", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	
	response = httptest.NewRecorder()
	healthz(response, request)
	
	if response.Code != 500 {
		t.Errorf("Expected status code 500, got %d", response.Code)
	}
	
	// Restore the original state
	healthy = originalHealthy
}
