package appcontext

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestClientID(t *testing.T) {
	testID := uuid.New()
	tests := []struct {
		name           string
		ctx            context.Context
		expectedClient uuid.UUID
	}{
		{
			name:           "ClientID does not exist in context",
			ctx:            context.Background(),
			expectedClient: uuid.Nil,
		},
		{
			name:           "ClientID exists in context",
			ctx:            context.WithValue(context.Background(), KeyClientID, testID.String()),
			expectedClient: testID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientID := ClientID(tt.ctx)
			if clientID != tt.expectedClient {
				t.Errorf("got %v, want %v", clientID, tt.expectedClient)
			}
		})
	}
}
