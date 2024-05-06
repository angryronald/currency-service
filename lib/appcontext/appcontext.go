package appcontext

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	// KeyClientID represents the Current Client in http server context
	KeyClientID contextKey = "ClientID"
)

// ClientID gets current client from the context
func ClientID(ctx context.Context) uuid.UUID {
	currentClientAccess := (ctx).Value(KeyClientID)
	if currentClientAccess != nil {
		v := uuid.MustParse(currentClientAccess.(string))
		return v
	}
	return uuid.Nil
}
