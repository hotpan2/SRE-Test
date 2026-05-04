package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	healthzHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "OK" {
		t.Errorf("expected body OK, got %s", w.Body.String())
	}
}

func TestGreetingHandlerWithEnv(t *testing.T) {
	os.Setenv("GO_SERVICE_NAME", "test-service")
	defer os.Unsetenv("GO_SERVICE_NAME")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	greetingHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "test-service") {
		t.Errorf("expected greeting with test-service, got %s", w.Body.String())
	}
}

func TestGreetingHandlerDefaultName(t *testing.T) {
	os.Unsetenv("GO_SERVICE_NAME")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	greetingHandler(w, req)

	if !strings.Contains(w.Body.String(), "go-service") {
		t.Errorf("expected default name go-service, got %s", w.Body.String())
	}
}

func TestNewMuxRoutes(t *testing.T) {
	mux := newMux()

	tests := []struct {
		path           string
		expectedStatus int
	}{
		{"/healthz", http.StatusOK},
		{"/", http.StatusOK},
		{"/metrics", http.StatusOK},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != tt.expectedStatus {
			t.Errorf("path %s: expected %d, got %d", tt.path, tt.expectedStatus, w.Code)
		}
	}
}
