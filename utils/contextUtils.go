package utils

import (
	"TODO-LIST/commons"
	"context"
	"fmt"
	"net/http"
)

type contextKey string

const clientIDKey = contextKey("client_id")

// AddClientIDToContext adds the client ID to the request context
func AddClientIDToContext(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, clientIDKey, clientID)
}

// GetClientIDFromContext retrieves the client ID from the request context
func GetClientIDFromContext(ctx context.Context) (string, error) {
	clientID, ok := ctx.Value(clientIDKey).(string)
	if !ok {
		return "", fmt.Errorf("client ID not found or invalid")
	}
	return clientID, nil
}

// GetClientCredentials retrieves the ClientCredentials from the context of the request.
func GetClientCredentials(r *http.Request) (*commons.ClientCredentials, bool) {
	creds, ok := r.Context().Value("creds").(*commons.ClientCredentials)
	return creds, ok
}
