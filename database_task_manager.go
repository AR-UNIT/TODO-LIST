package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq" // Import the pq driver for database/sql
)

// DBConfig holds the database connection parameters
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SSLMode  string
}

// DatabaseTaskManager manages tasks in a PostgreSQL database
type DatabaseTaskManager struct {
	db *sql.DB
	mu sync.Mutex
}

func (dtm *DatabaseTaskManager) LoadTasks() ([]Task, int, error) {
	//TODO implement me
	panic("implement me")
}

func (dtm *DatabaseTaskManager) SaveTasks() error {
	//TODO implement me
	panic("implement me")
}

// Initialize establishes the database connection
func (dtm *DatabaseTaskManager) Initialize() {
	// Define the database configuration
	config := DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
		DbName:   "postgres",
		SSLMode:  "disable",
	}

	// Initialize the database connection
	db, err := InitializeDB(config)
	if err != nil || db == nil {
		log.Println("Error initializing database:", err)
		return
	}

	// Assign the established connection to the struct
	dtm.db = db
	log.Println("Database connection initialized successfully.")
}

// InitializeDB initializes and returns a database connection
func InitializeDB(config DBConfig) (*sql.DB, error) {
	// Build the connection string for lib/pq
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
		config.SSLMode,
	)

	// Open the database connection
	db, err := sql.Open("postgres", dsn) // Use "postgres" for lib/pq
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database!")
	return db, nil
}

// AddTask inserts a new task into the database
func (dtm *DatabaseTaskManager) AddTask(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO TODO.tasks (description, completed) VALUES ($1, $2) RETURNING id`
	err := dtm.db.QueryRow(query, task.Description, false).Scan(&task.ID)
	if err != nil {
		http.Error(w, "Error inserting task into database", http.StatusInternalServerError)
		return
	}

	task.Completed = false
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (dtm *DatabaseTaskManager) ListTasks(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	rows, err := dtm.db.Query(`SELECT id, description, completed FROM TODO.tasks`)
	if err != nil {
		log.Println("Error querying tasks:", err) // Log the actual error
		http.Error(w, "Error retrieving tasks from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Completed); err != nil {
			log.Println("Error scanning task:", err) // Log scanning errors
			http.Error(w, "Error reading tasks from database", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CompleteTask marks a task as completed in the database
func (dtm *DatabaseTaskManager) CompleteTask(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	id := r.URL.Query().Get("id")
	query := `UPDATE TODO.tasks SET completed = true WHERE id = $1`
	res, err := dtm.db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task marked as completed")
}

// DeleteTask deletes a task from the database
func (dtm *DatabaseTaskManager) DeleteTask(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	id := r.URL.Query().Get("id")
	query := `DELETE FROM TODO.tasks WHERE id = $1`
	res, err := dtm.db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task deleted successfully")
}

// LazySave is a placeholder for additional database operations
func (dtm *DatabaseTaskManager) LazySave() {
	// This function can be used for periodic save operations if needed
}
