package main

// TaskManager defines the methods for managing tasks
type TaskManager interface {
	LoadTasks() error
	SaveTasks() error
	AddTask(description string)
	ListTasks()
	CompleteTask(id int)
	DeleteTask(id int)
}
