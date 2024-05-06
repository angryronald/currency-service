package http

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func TestRouteExists(t *testing.T) {
	// Create a chi.Mux router
	router := chi.NewRouter()

	// Register some routes
	router.Get("/authentications/otp/{phone}", func(w http.ResponseWriter, r *http.Request) {})
	router.Post("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/products/{productID}", func(w http.ResponseWriter, r *http.Request) {})

	// Test if the routes exist
	if !IsRouteExists(router, "/authentications/otp/{phone}") {
		t.Errorf("Expected route '/authentications/otp/{phone}' not found in the router")
	}

	if !IsRouteExists(router, "/users/{userID}") {
		t.Errorf("Expected route '/users/{userID}' not found in the router")
	}

	if !IsRouteExists(router, "/products/{productID}") {
		t.Errorf("Expected route '/products/{productID}' not found in the router")
	}

	if IsRouteExists(router, "/nonexistent/route") {
		t.Errorf("Unexpected route '/nonexistent/route' found in the router")
	}
}
