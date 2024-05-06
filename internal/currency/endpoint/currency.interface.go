package endpoint

import (
	"github.com/go-chi/chi"
)

type CurrencyEndpointInterface interface {
	RegisterRoute(r *chi.Mux) *chi.Mux
}
