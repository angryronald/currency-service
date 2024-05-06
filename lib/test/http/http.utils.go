package http

import (
	"strings"

	"github.com/go-chi/chi"
)

// Helper function to check if a route exists in the chi router
func IsRouteExists(router *chi.Mux, routePattern string) bool {
	// Iterate through the registered routes and check if the pattern matches
	for _, route := range router.Routes() {
		if strings.Contains(route.Pattern, routePattern) {
			return true
		}
		if route.SubRoutes != nil {
			for _, subroute := range route.SubRoutes.Routes() {
				if strings.Contains(subroute.Pattern, routePattern) {
					return true
				}
			}
		}
	}
	return false
}
