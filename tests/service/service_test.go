package service_test

import (
	"testing"

	"game-platform/internal/service"
)

// PortalAPI and S3Client tests are integration tests requiring real services.
// They are skipped here; unit tests for handlers/services use mocks instead.

func TestPortalAPI_Construct(t *testing.T) {
	// Verify PortalAPI can be constructed without panicking
	portalAPI := service.NewPortalAPI("", "")
	if portalAPI == nil {
		t.Fatal("Expected PortalAPI, got nil")
	}
	if portalAPI.BaseURL() != "" {
		t.Errorf("Expected empty base URL, got %s", portalAPI.BaseURL())
	}
}

func TestPortalAPI_ConstructWithTimeout(t *testing.T) {
	portalAPI := service.NewPortalAPIWithTimeout("http://localhost:8080", "test-key", 5)
	if portalAPI == nil {
		t.Fatal("Expected PortalAPI, got nil")
	}
	if portalAPI.BaseURL() != "http://localhost:8080" {
		t.Errorf("Expected http://localhost:8080, got %s", portalAPI.BaseURL())
	}
}
