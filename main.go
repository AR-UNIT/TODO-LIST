package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"sync"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks []Task
var taskID int
var mu sync.Mutex

func errorHandler(err error, errorType string) {
	if err != nil {
		fmt.Println(errorType, err)
		return
	}
}

func main() {
	// TODO:
	// hardcoded storage type, to get the task_manager for different storage types
	//taskStorageType := "file"
	taskStorageType := "database"
	manager, err := GetTaskManager(taskStorageType)
	errorHandler(err, ERROR_CREATING_TASK_MANAGER)

	manager.Initialize()
	errorHandler(err, ERROR_LOADING_TASKS)

	// Define the routes and handlers
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":

			manager.ListTasks(w, r)
		case "POST":
			manager.AddTask(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/complete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			manager.CompleteTask(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			manager.DeleteTask(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	// Save tasks to the file when the application exits
	defer manager.LazySave()
}
