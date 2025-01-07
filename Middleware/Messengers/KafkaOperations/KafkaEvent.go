package KafkaOperations

import (
	"encoding/json"
	"github.com/google/uuid" // Import UUID generation library
	"log"
	"net/http"
	"time"
)

// KafkaEvent represents a single Kafka event message
type KafkaEvent struct {
	EventID     string            `json:"event_id"` // Unique ID for the event
	EventType   string            `json:"event_type"`
	Timestamp   string            `json:"timestamp"` // ISO 8601 formatted timestamp
	Endpoint    string            `json:"endpoint"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	Payload     interface{}       `json:"payload"`
	ClientID    string            `json:"user_id"` // Optional if authenticated user
}

// NewKafkaEvent creates a new KafkaEvent with the current timestamp
func NewKafkaEvent(eventType string, endpoint string, headers map[string]string, queryParams map[string]string, payload interface{}, clientID string) KafkaEvent {
	return KafkaEvent{
		EventID:     generateUniqueID(), // Generate a unique ID for the event
		EventType:   eventType,
		Timestamp:   time.Now().Format(time.RFC3339), // ISO 8601 format
		Endpoint:    endpoint,
		Headers:     headers,
		QueryParams: queryParams,
		Payload:     payload,
		ClientID:    clientID,
	}
}

// Helper function to generate a unique event ID
func generateUniqueID() string {
	return uuid.New().String() // Generates a new UUID
}

func extractHeaders(r *http.Request) map[string]string {
	headers := map[string]string{}
	for key, values := range r.Header {
		headers[key] = values[0] // Use the first value if multiple
	}
	return headers
}

func extractQueryParams(r *http.Request) map[string]string {
	queryParams := map[string]string{}
	for key, values := range r.URL.Query() {
		queryParams[key] = values[0] // Use the first value if multiple
	}
	return queryParams
}

func extractPayload(r *http.Request) map[string]interface{} {
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error decoding request body:", err)
	}
	return payload
}
