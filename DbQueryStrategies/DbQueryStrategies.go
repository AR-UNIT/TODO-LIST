package DbQueryStrategies

import (
	"TODO-LIST/commons"
	_ "github.com/lib/pq" // Import the pq driver for database/sql
)

// DatabaseQueryStrategy defines an interface for database operations
type DatabaseQueryStrategy interface {
	AddTask(task commons.Task) (int, error)
	ListTasks() ([]commons.Task, error)
	CompleteTask(id int) (int64, error) // returns the number of rows affected, or error
	DeleteTask(id int) (int64, error)   // returns the number of rows affected, or error
}
