package constants

const (
	ERROR_LOADING_TASKS         = "Error loading tasks"
	ERROR_CREATING_TASK_MANAGER = "Error creating task manager"
)

const (
	CREATE_TASK   = "CreateTask"
	COMPLETE_TASK = "CompleteTask"
	DELETE_TASK   = "DeleteTask"
	LIST_TASK     = "ListTask"
)

// Error Messages
const (
	ErrorInvalidMethod          = "Invalid request method"
	ErrorInvalidInput           = "Invalid input"
	ErrorUnauthorized           = "Unauthorized access"
	ErrorDatabaseInitialization = "Error initializing database connection"
	ErrorGeneratingToken        = "Error generating token"
	ErrorInvalidCredentials     = "Invalid credentials"
	ErrorTaskNotFound           = "Task not found"
)

// Success Messages
const (
	SuccessTaskCompleted = "Task marked as complete"
	SuccessTaskDeleted   = "Task successfully deleted"
	SuccessTaskAdded     = "Task successfully added"
)
