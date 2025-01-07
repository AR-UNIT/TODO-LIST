package auth

import (
	"TODO-LIST/Middleware/Authenticators/jwt"
	dbTaskManager "TODO-LIST/TaskManagers"
	"TODO-LIST/commons"
	"TODO-LIST/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Assuming you're using PostgreSQL
	"log"
	"net/http"
	"os"
	"strconv"
)

// USED BY CLIENT AT ENDPOINT /api/authenticate TO GET THE JWT TOKEN

// ClientCredentials represents the payload for login requests

// AuthenticateAndProvideJWT handles client login and JWT generation
func AuthenticateAndProvideJWT(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthenticateAndProvideJWT")

	// rate limiter would have been called first, which already puts the client credentials in ctx
	creds, ok := utils.GetClientCredentials(r)
	if !ok {
		http.Error(w, "Failed to retrieve client credentials from payload", http.StatusInternalServerError)
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

	fmt.Println("Database connection established, going to do the client lookup now")
	var clientLookup commons.ClientCredentials
	query := "SELECT client_id, client_secret FROM TODO.clients WHERE client_id = $1"
	err = db.QueryRow(query, creds.ClientID).Scan(&clientLookup.ClientID, &clientLookup.ClientSecret)
	if err != nil || clientLookup.ClientSecret != creds.ClientSecret {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Println("Successfully authenticated user: ", clientLookup.ClientID, clientLookup.ClientSecret)

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
	fmt.Println("generated jwt token:  ", token)
}

// Initialize sets up to the DB connection storing client details and returns it
func Initialize() (*sql.DB, error) {
	// Set up DB connection (you might want to use a connection pool in real applications)
	// Define the database configuration

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming environment variables are set")
	}
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	// Convert DB_PORT from string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	// Define the database configuration
	config := commons.DBConfig{
		Host:     host,
		Port:     port, // Use the integer value
		User:     user,
		Password: password,
		DbName:   dbName,
		SSLMode:  sslMode,
	}

	// Initialize the database connection
	db, err := dbTaskManager.InitializeDB(config)
	fmt.Println("returned back to initialize in jwt handler after doing db initialization")

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
