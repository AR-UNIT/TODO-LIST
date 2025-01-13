package main

import (
	handlerJWT "TODO-LIST/Handlers/auth"
	authJWT "TODO-LIST/Middleware/Authenticators/jwt"
	kafkaOperations "TODO-LIST/Middleware/Messengers/KafkaOperations"
	"TODO-LIST/Middleware/MetricsCollectors/Prometheus"
	"TODO-LIST/Middleware/RateLimiters"
	"TODO-LIST/TaskManagers"
	redisCache "TODO-LIST/caches/Redis"
	"TODO-LIST/constants"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"strconv"
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
	err := godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err) // Log the error with more details
	}

	taskStorageType := os.Getenv("TASK_STORAGE_TYPE")
	taskManager, err := TaskManagers.GetTaskManager(taskStorageType)
	errorHandler(err, constants.ERROR_CREATING_TASK_MANAGER)
	taskManager.Initialize()
	errorHandler(err, constants.ERROR_LOADING_TASKS)

	kafkaBrokerAddress := os.Getenv("KAFKA_BROKER_ADDRESS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaConsumerGroupId := os.Getenv("KAFKA_CONSUMER_GROUP_ID")

	err = kafkaOperations.InitKafkaProducer(kafkaBrokerAddress, kafkaTopic)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	defer kafkaOperations.CloseProducer() // Close Kafka producer on exit

	// Define the Kafka consumer configuration
	kafkaConfig := kafkaOperations.KafkaConsumerConfig{
		BrokerAddress: kafkaBrokerAddress,   // Kafka broker address
		Topic:         kafkaTopic,           // Kafka topic for task-related events
		GroupID:       kafkaConsumerGroupId, // Consumer group ID
		TaskManager:   taskManager,          // Pass the initialized taskManager instance
	}

	// Start the Kafka consumer in a separate goroutine
	go kafkaOperations.StartKafkaConsumer(kafkaConfig)

	// Continue with other application setup (e.g., HTTP server, etc.)
	log.Println("Kafka consumer started...")

	rateLimiterLimit := os.Getenv("RATE_LIMITER_LIMIT")
	rateLimiterBust := os.Getenv("RATE_LIMITER_BURST")

	rateLimit, err := strconv.Atoi(rateLimiterLimit)
	if err != nil {
		log.Fatalf("Invalid rateLimit value: %v", err)
	}

	rateBurst, err := strconv.Atoi(rateLimiterBust)
	if err != nil {
		log.Fatalf("Invalid rateBurst value: %v", err)
	}
	/*
		The initial burst allows 10 requests to be sent quickly (all at once or within a fraction of a second).
		Once that burst is consumed, the system switches to enforcing the 5 requests per second rate limit.
		If the client exceeds the 5 requests per second, further requests will be rejected until the next second when the limit resets.

		burst is replenished based on rate specified
	*/
	rateLimiter := RateLimiters.NewRateLimiter(rate.Limit(rateLimit), rateBurst)

	redisCache.InitRedis()

	go Prometheus.ExposeMetrics()

	/* 	Register the authentication endpoint
	Handles client login and JWT generation
	ratelimiter is called first, and rate is applied after identifying client
	*/
	http.Handle("/api/authenticate", rateLimiter.Apply(http.HandlerFunc(handlerJWT.AuthenticateAndProvideJWT)))

	/* registering endpoints for crud operations */
	http.Handle("/tasks/complete",
		Prometheus.CountRequests( // Apply Prometheus countRequests middleware
			rateLimiter.Apply( // Apply the rate limiter
				authJWT.AuthenticateJWT( // Apply JWT authentication
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						log.Printf("Received request for /tasks/complete with method: %s", r.Method)
						if r.Method == constants.HTTPMethodPatch {
							kafkaOperations.TaskHandler(constants.COMPLETE_TASK, w, r)
						} else {
							http.Error(w, constants.ErrorInvalidMethod, constants.StatusMethodNotAllowed)
						}
					}),
				),
			),
		),
	)

	http.Handle("/tasks/delete",
		Prometheus.CountRequests( // Apply Prometheus countRequests middleware
			rateLimiter.Apply( // Apply the rate limiter
				authJWT.AuthenticateJWT( // Apply JWT authentication
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						log.Printf("Received request for /tasks/delete with method: %s", r.Method)
						if r.Method == constants.HTTPMethodDelete {
							kafkaOperations.TaskHandler(constants.DELETE_TASK, w, r)
						} else {
							http.Error(w, constants.ErrorInvalidMethod, constants.StatusMethodNotAllowed)
						}
					}),
				),
			),
		),
	)
	http.Handle("/tasks",
		Prometheus.CountRequests( // Apply Prometheus countRequests middleware
			rateLimiter.Apply( // Apply the rate limiter
				authJWT.AuthenticateJWT( // Apply JWT authentication
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						log.Printf("Received request for /tasks with method: %s", r.Method)

						switch r.Method {
						case constants.HTTPMethodGet:
							taskID := r.URL.Query().Get("id")

							// Check Redis cache first
							task, err := redisCache.GetTaskFromCache(taskID)
							if err != nil {
								fmt.Println("Error while making cache call for tasks")
							}
							if task != nil {
								// Cache hit, return the task directly
								w.WriteHeader(http.StatusOK)
								json.NewEncoder(w).Encode(task)
							} else {
								// Cache miss, fetch from DB
								taskManager.ListTasks(w, r)
							}

						case constants.HTTPMethodPost:
							kafkaOperations.TaskHandler(constants.CREATE_TASK, w, r)
						default:
							http.Error(w, constants.ErrorInvalidMethod, constants.StatusMethodNotAllowed)
						}
					}),
				),
			),
		),
	)

	// Start the server
	apiPort := os.Getenv("API_PORT")
	fmt.Printf("Server started on http://localhost%s\n", apiPort)
	if err := http.ListenAndServe(apiPort, nil); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	// TODO:
	// THIS WILL NEVER BE CALLED WITHOUT A GRACEFUL SHUTDOWN HANDLER,
	// ALSO THIS IS NOT USEFUL IF WE ARE NOT USING DOING PERSISTANCE OF CACHED UPDATED TO DB BEFORE UPON FAILURE,
	// THIS IS GENERALLY NOT FAULT TOLERANT WITHOUT RECOVERY ROUTINE, NEED TO FIX
	// ALL OPERATIONS ARE ALWAYS BEEN DONE DIRECTLY TO THE DB
	// Save tasks to the file when the application exits
	defer taskManager.LazySave()
}
