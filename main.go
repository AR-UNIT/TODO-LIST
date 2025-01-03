package main

import (
	handlerJWT "TODO-LIST/Handlers/auth"
	authJWT "TODO-LIST/Middleware/Authenticators/jwt" // Authenticator package
	"TODO-LIST/Middleware/RateLimiters"
	"TODO-LIST/TaskManagers"
	"TODO-LIST/constants"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
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
	errorHandler(err, constants.ERROR_CREATING_TASK_MANAGER)
	manager.Initialize()
	errorHandler(err, constants.ERROR_LOADING_TASKS)

	/*
		The initial burst allows 10 requests to be sent quickly (all at once or within a fraction of a second).
		Once that burst is consumed, the system switches to enforcing the 5 requests per second rate limit.
		If the client exceeds the 5 requests per second, further requests will be rejected until the next second when the limit resets.

		burst is replenished based on rate specified
	*/
	rateLimiter := RateLimiters.NewRateLimiter(rate.Limit(5), 10)

	// Register the authentication endpoint
	// Handles client login and JWT generation
	http.Handle("/api/authenticate", rateLimiter.Apply(http.HandlerFunc(handlerJWT.AuthenticateClient)))

	http.Handle("/tasks", rateLimiter.Apply(authJWT.AuthenticateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case constants.HTTPMethodGet:
			manager.ListTasks(w, r)
		case constants.HTTPMethodPost:
			manager.AddTask(w, r)
		case constants.HTTPMethodPatch:
			manager.CompleteTask(w, r)
		case constants.HTTPMethodDelete:
			manager.DeleteTask(w, r)
		default:
			http.Error(w, constants.ErrorInvalidMethod, constants.StatusMethodNotAllowed)
		}
	}))))

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
