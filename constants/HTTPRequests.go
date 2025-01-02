package constants

// HTTP Method Constants
const (
	HTTPMethodGet    = "GET"
	HTTPMethodPost   = "POST"
	HTTPMethodPatch  = "PATCH"
	HTTPMethodDelete = "DELETE"
)

// Database Queries
const (
	QueryClientLookup = "SELECT client_id, client_secret FROM clients WHERE client_id = $1"
	// Add more query strings as needed
)

// HTTP Status Codes (for readability)
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusMethodNotAllowed    = 405
	StatusInternalServerError = 500
)
