package leadership

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakater/Reloader/internal/pkg/testutil"
	"github.com/stakater/Reloader/pkg/kube"
	"github.com/stretchr/testify/assert"
)

func TestLeadership(t *testing.T) {
	// Setup for this specific test
	testutil.SetupTestEnvironment()

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
		testutil.TeardownTestEnvironment()
	}()

	// Test NewLeadership function
	t.Run("TestNewLeadership", func(t *testing.T) {
		client, _, _ := testutil.CreateFakeClientSet()
		leadership := NewLeadership(client)
		assert.NotNil(t, leadership, "Leadership should not be nil")
		assert.Equal(t, client, leadership.client, "Client should be set correctly")
	})

	// Test healthz endpoint
	t.Run("TestHealthzEndpoint", func(t *testing.T) {
		// Reset the healthy state to true
		m.Lock()
		healthy = true
		m.Unlock()

		// Create a request to pass to our handler
		req, err := http.NewRequest("GET", "/live", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthz)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder
		handler.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		// Check the response body
		assert.Equal(t, "alive", rr.Body.String(), "handler returned unexpected body")
	})
}
