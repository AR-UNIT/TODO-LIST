package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// FileLoader struct for loading tasks from a file
type FileTaskManager struct {
	FilePath string
}

func (fl *FileTaskManager) Initialize() {

}

func (fl *FileTaskManager) SaveTasks() error {
	//TODO implement me
	panic("implement me")
}

func (fl *FileTaskManager) AddTask(w http.ResponseWriter, r *http.Request) {
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

func (fl *FileTaskManager) ListTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if len(tasks) == 0 {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (fl *FileTaskManager) CompleteTask(w http.ResponseWriter, r *http.Request) {
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

func (fl *FileTaskManager) DeleteTask(w http.ResponseWriter, r *http.Request) {
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

// LoadTasks method for FileLoader
func (fl *FileTaskManager) LoadTasks() ([]Task, int, error) {
	file, err := os.Open(fl.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, 0, nil // Return empty list and taskID 0 if file doesn't exist
		}
		return nil, 0, fmt.Errorf("error loading tasks from file: %v", err)
	}
	defer file.Close()

	var loadedTasks []Task
	if err := json.NewDecoder(file).Decode(&loadedTasks); err != nil {
		return nil, 0, fmt.Errorf("error decoding tasks: %v", err)
	}

	// Determine the next task ID to avoid duplicate IDs
	maxID := 0
	for _, task := range loadedTasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return loadedTasks, maxID, nil
}

// Save tasks to a file through lazy calls to save, need to figure thist out, and also
// how to maintain cache of most recent data in memory
func (fl *FileTaskManager) LazySave() {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		os.Exit(1)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(tasks)
}
