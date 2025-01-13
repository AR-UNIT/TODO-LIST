package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var RedisClient *redis.Client

func InitRedis() {
	// Get Redis host and port from environment variables
	redisHost := os.Getenv("REDIS_HOST") // Should be 'redis' as defined in docker-compose.yml
	redisPort := os.Getenv("REDIS_PORT") // Default port is 6379

	if redisHost == "" || redisPort == "" {
		redisHost = "localhost" // Default to localhost if not set
		redisPort = "6379"      // Default port
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Initialize Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0, // Default DB is 0
	})

	// Test the connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}

	fmt.Println("Successfully connected to Redis")
}

// CacheTask will cache a task in Redis
func CacheTask(taskID string, task interface{}) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("could not marshal task: %v", err)
	}

	// Set the task in Redis with a 1-hour expiration time
	err = RedisClient.Set(context.Background(), taskID, taskJSON, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("could not cache task: %v", err)
	}
	return nil
}

// GetTaskFromCache retrieves a task from Redis cache
func GetTaskFromCache(taskID string) (interface{}, error) {
	taskJSON, err := RedisClient.Get(context.Background(), taskID).Result()
	if err == redis.Nil {
		// Cache miss, return nil
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not retrieve task from cache: %v", err)
	}

	var task interface{}
	err = json.Unmarshal([]byte(taskJSON), &task)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal task: %v", err)
	}
	fmt.Println("fetched task with id: ", taskID, task)
	return task, nil
}

// DeleteTaskFromCache invalidates a cached task
func DeleteTaskFromCache(taskID string) error {
	err := RedisClient.Del(context.Background(), taskID).Err()
	if err != nil {
		return fmt.Errorf("could not delete task from cache: %v", err)
	}
	return nil
}

// CacheTaskList caches a list of tasks (for GET requests)
func CacheTaskList(tasks []interface{}) error {
	// Serialize the list of tasks
	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("could not marshal tasks list: %v", err)
	}

	// Cache with a key like "tasks_list" for all tasks
	err = RedisClient.Set(context.Background(), "tasks_list", tasksJSON, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("could not cache tasks list: %v", err)
	}
	return nil
}

// GetTaskListFromCache retrieves the list of tasks from cache
func GetTaskListFromCache() ([]interface{}, error) {
	tasksJSON, err := RedisClient.Get(context.Background(), "tasks_list").Result()
	if err == redis.Nil {
		// Cache miss, return nil
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not retrieve tasks list from cache: %v", err)
	}

	var tasks []interface{}
	err = json.Unmarshal([]byte(tasksJSON), &tasks)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal tasks list: %v", err)
	}
	return tasks, nil
}
