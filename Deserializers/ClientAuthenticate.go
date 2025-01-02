package Deserializers

import (
	"TODO-LIST/commons"
	"context"
	"net/http"
)

// ClientAuthenticate extracts the client credentials from request headers
// and attaches them to the request context.
func ClientAuthenticate(w http.ResponseWriter, r *http.Request) (*commons.ClientCredentials, *http.Request) {
	// Extract the client credentials from headers
	clientName := r.Header.Get("client_name")
	clientID := r.Header.Get("client_id")
	clientSecret := r.Header.Get("client_secret")

	// If either clientID or clientSecret is missing, return an error
	if clientID == "" || clientSecret == "" {
		http.Error(w, "Client ID and Client Secret required", http.StatusBadRequest)
		return nil, r
	}

	// Create the ClientCredentials struct
	creds := &commons.ClientCredentials{
		ClientName:   clientName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Attach creds to the request context
	ctx := context.WithValue(r.Context(), "creds", creds)
	r = r.WithContext(ctx) // Return the modified request with the new context

	// Return a pointer to creds and the updated request
	return creds, r
}
