package KafkaOperations

import (
	"TODO-LIST/Deserializers"
	"TODO-LIST/TaskManagers"
	"TODO-LIST/constants"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

// KafkaConsumerConfig holds the configuration for the Kafka consumer
type KafkaConsumerConfig struct {
	BrokerAddress string
	Topic         string
	GroupID       string
	TaskManager   TaskManagers.TaskManager // Include TaskManager instance
}

// StartKafkaConsumer initializes and starts the Kafka consumer
func StartKafkaConsumer(config KafkaConsumerConfig) {
	// Create a new Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.BrokerAddress},
		Topic:    config.Topic,
		GroupID:  config.GroupID,
		MaxBytes: 10e6, // 10MB max per message
	})

	log.Printf("Kafka consumer started for topic %s on broker %s", config.Topic, config.BrokerAddress)

	for {
		// Read message from Kafka
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		log.Printf("Message received: key=%s, value=%s", string(msg.Key), string(msg.Value))

		// Deserialize the Kafka message into a KafkaEvent
		var event KafkaEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Error deserializing Kafka message:", err)
			continue
		}

		// Process the Kafka event
		handleKafkaEvent(event, config.TaskManager)
	}
}

// handleKafkaEvent processes a Kafka event
func handleKafkaEvent(event KafkaEvent, taskManager TaskManagers.TaskManager) {
	// Add your custom event processing logic here
	log.Printf("Processing event: ID=%s, Type=%s, Endpoint=%s, ClientID=%s",
		event.EventID, event.EventType, event.Endpoint, event.ClientID)

	if taskManager == nil {
		log.Fatalf("Failed to create TaskManager: %v")
	}
	log.Println(taskManager)

	log.Printf("Raw payload before deserialization: %s", event.Payload)
	taskInput := Deserializers.DeserializeTaskInput(event.Payload)
	if taskInput == nil {
		log.Println("Error: Failed to deserialize task input")
		return
	}

	switch event.EventType {

	case constants.CREATE_TASK:
		taskManager.AddTask(taskInput)

	case constants.COMPLETE_TASK:
		log.Println("Handle CompleteTask event")
		// Use taskManager to complete a task
		fmt.Println("at complete task")

		//taskManager.CompleteTask(event.Payload)
		//if err != nil {
		//	log.Printf("Error completing task: %v", err)
		//}

	case constants.DELETE_TASK:
		log.Println("Handle DeleteTask event")
		// Use taskManager to delete a task
		//err := taskManager.DeleteTask(event.Payload)
		//if err != nil {
		//	log.Printf("Error deleting task: %v", err)
		//}

	default:
		log.Printf("Unhandled event type: %s", event.EventType)
	}
}
