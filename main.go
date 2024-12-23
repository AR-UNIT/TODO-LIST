package main

import (
	"TODO-LIST/TaskManagers"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func errorHandler(err error, errorType string) {
	if err != nil {
		fmt.Println(errorType, err)
		return
	}
}

func main() {
	// TODO:
	// hardcoded storage type, to get the task_manager for different storage types, make dynamic
	taskStorageType := "postgresDb"
	manager, err := TaskManagers.GetTaskManager(taskStorageType)
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

	// TODO:
	// THIS WILL NEVER BE CALLED WITHOUT A GRACEFUL SHUTDOWN HANDLER,
	// ALSO THIS IS NOT USEFUL IF WE ARE NOT USING A CACHE TO STORE OPERATIONS IN MEMORY,
	// ALL OPERATIONS ARE ALWAYS BEEN DONE DIRECTLY TO THE DB
	// Save tasks to the file when the application exits
	defer manager.LazySave()
}
