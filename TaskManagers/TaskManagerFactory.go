package TaskManagers

import (
	"TODO-LIST/commons"
	"fmt"
	"net/http"
)

// TaskLoader interface for task loading methods
type TaskManager interface {
	Initialize()
	SaveTasks() error
	AddTask(task *commons.TaskInputModel)
	ListTasks(w http.ResponseWriter, r *http.Request)
	CompleteTask(taskId string)
	DeleteTask(taskId string)
	LazySave()
}

// Factory method to create the appropriate TaskLoader
func GetTaskManager(sourceType string) (TaskManager, error) {
	switch sourceType {
	case "postgresDb":
		return &DatabaseTaskManager{}, nil
	default:
		return nil, fmt.Errorf("unsupported task loader type: %s", sourceType)
	}
}
