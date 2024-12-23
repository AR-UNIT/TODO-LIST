package TaskManagers

import (
	"TODO-LIST/DbQueryStrategies"
	"TODO-LIST/commons"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // Import the pq driver for database/sql
	"log"
	"net/http"
	"strconv"
	"sync"
)

// SelectStrategy determines which strategy to use based on a given mode
func SelectStrategy(db *sql.DB, queryStrategy string) DbQueryStrategies.DatabaseQueryStrategy {
	switch queryStrategy {
	case "rowLockingStrategy":
		return &DbQueryStrategies.PostgresRowLockingStrategy{
			BasePostgresStrategy: DbQueryStrategies.BasePostgresStrategy{Db: db},
		}
	default:
		return &DbQueryStrategies.PostgresStrategy{
			BasePostgresStrategy: DbQueryStrategies.BasePostgresStrategy{Db: db},
		}
	}
}

func (dtm *DatabaseTaskManager) SaveTasks() error {
	//TODO implement me
	panic("implement me")
}

// InitializeDB initializes and returns a database connection
func InitializeDB(config commons.DBConfig) (*sql.DB, error) {
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

// DatabaseTaskManager manages tasks using a database and a strategy
type DatabaseTaskManager struct {
	Db       *sql.DB
	strategy DbQueryStrategies.DatabaseQueryStrategy
	mu       sync.Mutex
}

func (dtm *DatabaseTaskManager) LazySave() {
	//TODO implement me
	panic("implement me")
}

func (dtm *DatabaseTaskManager) Initialize() {
	// Define the database configuration
	config := commons.DBConfig{
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
	dtm.Db = db
	log.Println("Database connection initialized successfully.")

	// Choose whether to use row locking (example: based on an environment variable or config)
	dbQueryStrategy := "rowLockingStrategy" // Set this based on your application's requirements

	// Select the strategy and assign it
	dtm.strategy = SelectStrategy(dtm.Db, dbQueryStrategy)
	log.Println("Current query strategy: ", dbQueryStrategy)
}

// AddTask adds a task to the database using the strategy
func (dtm *DatabaseTaskManager) AddTask(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	var task commons.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Add the task using the strategy method (no need to pass Db explicitly)
	id, err := dtm.strategy.AddTask(task)
	if err != nil {
		http.Error(w, "Error inserting task into database", http.StatusInternalServerError)
		return
	}

	task.ID = id
	task.Completed = false
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// ListTasks lists tasks from the database using the strategy
func (dtm *DatabaseTaskManager) ListTasks(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	// List the tasks using the strategy method (no need to pass Db explicitly)
	tasks, err := dtm.strategy.ListTasks()
	if err != nil {
		http.Error(w, "Error retrieving tasks from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CompleteTask marks a task as completed using the strategy
func (dtm *DatabaseTaskManager) CompleteTask(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	// Convert id from string to int
	id := r.URL.Query().Get("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Call CompleteTask on the strategy (pass the task ID as an int)
	rowsAffected, err := dtm.strategy.CompleteTask(taskID)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task marked as completed, %d rows affected", rowsAffected)
}

// DeleteTask deletes a task using the strategy and returns the number of rows affected
func (dtm *DatabaseTaskManager) DeleteTask(w http.ResponseWriter, r *http.Request) {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	id := r.URL.Query().Get("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Call DeleteTask on the strategy (pass the task ID as an int)
	rowsAffected, err := dtm.strategy.DeleteTask(taskID)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task deleted successfully, %d rows affected", rowsAffected)
}
