package main

import (
	"fmt"
	"net/http"
)

// TaskLoader interface for task loading methods
type TaskManager interface {
	LoadTasks() ([]Task, int, error)
	SaveTasks() error
	AddTask(w http.ResponseWriter, r *http.Request)
	ListTasks(w http.ResponseWriter, r *http.Request)
	CompleteTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	LazySave()
}

// Factory method to create the appropriate TaskLoader
func GetTaskManager(sourceType string) (TaskManager, error) {
	switch sourceType {
	case "file":
		return &FileTaskManager{FilePath: "tasks.json"}, nil
	//case "database":
	//	return &DatabaseTaskManager{}, nil
	default:
		return nil, fmt.Errorf("unsupported task loader type: %s", sourceType)
	}
}
