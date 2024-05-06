package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticateClient(t *testing.T) {
	// Define allowed clients for testing
	allowedClients := map[string]string{
		"token123": "client1",
		"token456": "client2",
	}

	// Create a request with authorization header
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("authorization", "Bearer token123")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create the middleware handler using the AuthenticateClient function
	handler := AuthenticateClient(allowedClients)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Serve HTTP using the middleware
	handler.ServeHTTP(rr, req)

	// Check if the request was handled properly
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test for unauthorized client
	reqUnauthorized := httptest.NewRequest("GET", "/", nil)
	reqUnauthorized.Header.Set("authorization", "Bearer invalidtoken")

	rrUnauthorized := httptest.NewRecorder()

	handler.ServeHTTP(rrUnauthorized, reqUnauthorized)

	if status := rrUnauthorized.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code for unauthorized client: got %v want %v", status, http.StatusForbidden)
	}
}
