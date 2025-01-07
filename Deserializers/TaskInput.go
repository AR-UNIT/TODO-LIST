package Deserializers

import (
	"TODO-LIST/commons"
	"encoding/json"
	"log"
)

func DeserializeTaskInput(payload interface{}) *commons.TaskInputModel {
	// Check if the payload is already a map[string]interface{}
	// This avoids unnecessary marshaling and unmarshaling
	if m, ok := payload.(map[string]interface{}); ok {
		// Directly deserialize from the map
		var taskInput commons.TaskInputModel
		err := json.Unmarshal(toBytes(m), &taskInput)
		if err != nil {
			log.Printf("Failed to deserialize payload: %v", err)
			return nil
		}
		return &taskInput
	}

	// If the payload is not a map, marshal it to JSON bytes
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload to JSON: %v", err)
		return nil
	}

	// Log the raw payload before deserialization
	log.Printf("Raw payload before deserialization: %s", string(rawPayload))

	// Deserialize the JSON payload into the TaskInputModel
	var taskInput commons.TaskInputModel
	err = json.Unmarshal(rawPayload, &taskInput)
	if err != nil {
		log.Printf("Failed to deserialize payload: %v", err)
		return nil
	}

	return &taskInput
}

// toBytes converts a map[string]interface{} to a byte slice for json.Unmarshal
func toBytes(m map[string]interface{}) []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println("Error marshalling map to bytes:", err)
		return nil
	}
	return bytes
}
