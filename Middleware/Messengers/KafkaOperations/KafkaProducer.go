package KafkaOperations

import (
	"TODO-LIST/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
)

var kafkaProducer *kafka.Writer

// InitKafkaProducer initializes the Kafka producer and returns an error if initialization fails
func InitKafkaProducer(brokerAddress, topic string) error {
	if brokerAddress == "" || topic == "" {
		return fmt.Errorf("broker address or topic cannot be empty")
	}

	kafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress), // Kafka broker address
		Topic:    topic,                    // Kafka topic
		Balancer: &kafka.LeastBytes{},      // Distribute messages based on size
	}

	// Check if the producer is successfully connected (try to write a dummy message to validate)
	err := kafkaProducer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
	})

	if err != nil {
		return fmt.Errorf("failed to initialize Kafka producer: %v", err)
	}

	log.Println("Kafka producer initialized successfully")
	return nil
}

// SendKafkaEvent sends an event to Kafka
func SendKafkaEvent(eventType string, headers map[string]string, queryParams map[string]string, payload interface{}, endpoint string, clientID string) {
	if kafkaProducer == nil {
		log.Println("Kafka producer is not initialized")
		return
	}

	// Build Kafka event
	event := NewKafkaEvent(eventType, endpoint, headers, queryParams, payload, clientID)

	// Serialize event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("Error serializing Kafka event:", err)
		return
	}

	// Send message to Kafka
	err = kafkaProducer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(event.EventID), // Use EventID as the message key
			Value: eventBytes,
		},
	)
	if err != nil {
		log.Println("Error sending Kafka message:", err)
	} else {
		log.Printf("Kafka event sent successfully: %v\n", event.EventID)
	}
}

// TaskHandler handles HTTP requests and sends Kafka events
func TaskHandler(eventType string, w http.ResponseWriter, r *http.Request) {
	// Extract relevant data
	headers := extractHeaders(r)
	queryParams := extractQueryParams(r)
	payload := extractPayload(r)
	clientID, err := utils.GetClientIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Client ID not found or invalid", http.StatusUnauthorized)
		return
	}
	fmt.Println("printing details extracted in task handler: ", headers, queryParams, payload, clientID)
	// Use the client ID for your logic
	fmt.Fprintf(w, "Client ID: %s", clientID)
	SendKafkaEvent(eventType, headers, queryParams, payload, r.URL.Path, clientID)

	// Handle HTTP response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task Created"))
}

// CloseProducer closes the Kafka producer and frees resources
func CloseProducer() {
	if kafkaProducer != nil {
		err := kafkaProducer.Close()
		if err != nil {
			log.Println("Error closing Kafka producer:", err)
		} else {
			log.Println("Kafka producer closed successfully")
		}
	} else {
		log.Println("Kafka producer is not initialized")
	}
}
