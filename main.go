package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Task struct to represent each task
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks []Task
var taskID int
var mu sync.Mutex

// Load tasks from a file
func loadTasksFromFile() {
	file, err := os.Open("tasks.json")
	if err != nil {
		// If the file doesn't exist, start with an empty task list
		if os.IsNotExist(err) {
			tasks = []Task{}
			taskID = 0
			return
		}
		fmt.Println("Error loading tasks:", err)
		os.Exit(1)
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&tasks)

	// Update taskID to avoid duplicate IDs
	for _, task := range tasks {
		if task.ID > taskID {
			taskID = task.ID
		}
	}
}

// Save tasks to a file
func saveTasksToFile() {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		os.Exit(1)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(tasks)
}

// Add a new task
func addTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	taskID++
	task.ID = taskID
	task.Completed = false
	tasks = append(tasks, task)

	// Respond with the added task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// List all tasks
func listTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if len(tasks) == 0 {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Complete a task
func completeTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// Delete a task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	// Load tasks from the file when the application starts
	loadTasksFromFile()

	// Define the routes and handlers
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			listTasks(w, r)
		case "POST":
			addTask(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/complete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			completeTask(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			deleteTask(w, r)
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
	defer saveTasksToFile()
}
