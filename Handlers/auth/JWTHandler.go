package auth

import (
	"TODO-LIST/Middleware/Authenticators/jwt"
	dbTaskManager "TODO-LIST/TaskManagers"
	"TODO-LIST/commons"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // Assuming you're using PostgreSQL
	"log"
	"net/http"
)

// ClientCredentials represents the payload for login requests

// AuthenticateClient handles client login and JWT generation
func AuthenticateClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthenticateClient")

	//creds, ok := r.Context().Value("creds").(commons.ClientCredentials)
	creds, ok := r.Context().Value("creds").(*commons.ClientCredentials)
	if !ok {
		http.Error(w, "Failed to retrieve client credentials", http.StatusInternalServerError)
		return
	}

	// Step 2: Now you have access to the creds, you can process the authentication
	fmt.Println("Received Client ID:", creds.ClientID)

	// Initialize DB connection
	db, err := Initialize()
	if err != nil {
		http.Error(w, "Error initializing database connection", http.StatusInternalServerError)
		return
	}
	defer db.Close() // Ensure the DB connection is closed after use

	// Validate credentials (replace with database lookup)
	var clientLookup commons.ClientCredentials
	query := "SELECT client_id, client_secret FROM TODO.clients WHERE client_id = $1"
	err = db.QueryRow(query, creds.ClientID).Scan(&clientLookup.ClientID, &clientLookup.ClientSecret)
	if err != nil || clientLookup.ClientSecret != creds.ClientSecret {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := jwt.GenerateJWT(clientLookup.ClientID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Respond with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Initialize sets up to the DB connection storing client details and returns it
func Initialize() (*sql.DB, error) {
	// Set up DB connection (you might want to use a connection pool in real applications)
	// Define the database configuration
	config := commons.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
		DbName:   "postgres",
		SSLMode:  "disable",
	}

	// Initialize the database connection
	db, err := dbTaskManager.InitializeDB(config)

	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	// Ensure the connection is alive
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	// Return the DB connection
	return db, nil
}
