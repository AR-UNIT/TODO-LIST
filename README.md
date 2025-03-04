# Distributed TODO API Project

## Project Overview

The project implements logically simple CRUD operations for a TODO application, while incorporating a wide range of tools and frameworks to enhance functionality and scalability.


## LLD Diagrams

### JWT Authentication Endpoint
![JWT Authentication Endpoint](./JWT%20TOKEN%20AUTH%20ENDPOINT.drawio.png)


### CRUD Endpoint Operations
![CRUD Endpoint Operations](./CRUD%20ENDPOINTS.drawio.png)


### HIGH LEVEL SOFTWARE ARCHITECTURE DIAGRAM
![high level design diagram](./high%20level%20design%20diagram.png)

## Technologies & Tools

- **PostgresDB**: Default storage method, with code extensible to different storage methods.
- **Docker & Kubernetes**: Application containerized and orchestrated using Docker and Kubernetes.
  - **Postgres Container**
  - **Zookeeper Container**
  - **Kafka Container**
  - **Redis Container**
  - **App Container**

## Middleware Features

- **JWT Authentication**: Ensures secure API calls with token-based authentication.
- **Rate Limiting**: Limits the number of requests per endpoint to prevent abuse.
- **Kafka**: Streams events from API requests and consumes events at the task manager for decoupled processing.
- **Redis**: Used for caching and performing batch updates to reduce DB calls as the application scales.

## Database Isolation Level and Concurrency

- **Row Locking For writes Strategy**: Prevents race conditions during write/update operations.
- **Read Committed Isolation Level**: Default isolation level in PostgresDB that provides certain guarantees but allows potential stale reads.

### Race Condition Safety in Read Committed and Row Locking

- **Read-Read**: Safe, no race condition(idempotent).
- **Read-Write**: Not safe; may lead to stale reads.
- **Write-Read**: Not safe; may lead to stale or inconsistent reads.
- **Write-Write**: Safe; row-level locking prevents conflicting writes.

## Architecture Diagram

### Distributed TODO API Architecture

1. **Clients**: 
   - Web or mobile apps communicate with the API over HTTPS.
2. **API Gateway**: 
   - Handles JWT authentication, rate limiting, and request routing.
   - handles api calls for task creation, update task to mark complete, delete task, and list tasks(in paginated + queried) format
3. **Backend Services**:
   - Written in Go implementing CRUD functionality for TODO items.
   - Follows factory and strategy patterns for extensibility.
     - Different query strategies can be implemented, and combined to form composite query strategy for the crud operations.
     - Task manager is provided via factory pattern, task manager is dependent on the mode of storage used for the tasks(PosgtesDB, etc).
   - Streams events to Kafka and integrates with Redis for caching.
4. **Kafka**: 
   - Streams API request events and handles event consumption.
5. **Redis**: 
   - Acts as a cache to improve query performance.
   - Performs batch updates to PostgreSQL for optimized database writes.
6. **PostgreSQL**: 
   - The primary database for storing TODO items.
7. **Deployment**:
   - Dockerized services orchestrated with Kubernetes.

USAGE (API CALLS):

Get jwt token:
$response = Invoke-WebRequest -Uri "http://localhost:8080/api/authenticate" -Method GET -Headers @{"client_id"="admin"; "client_secret"="password"; "Content-Type"="application/json"}
$token = ($response.Content | ConvertFrom-Json).token


GET TASKS / LIST TASKS
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method GET -Headers @{"Authorization"="Bearer $token";"client_id"="admin";"client_secret"="password";"Content-Type"="application/json"}

Invoke-WebRequest -Uri "http://localhost:8080/tasks?page=1&limit=10" -Method GET -Headers @{"Authorization"="Bearer $token";"client_id"="admin";"client_secret"="password";"Content-Type"="application/json"}


ADD TASKS 
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -Headers @{"Authorization"="Bearer $token"; "client_id"="admin"; "client_secret"="password";"Content-Type"="application/json"} -Body '{"description": "Buy groceries"}'

COMPLETE TASKS
Invoke-WebRequest -Uri http://localhost:8080/tasks/complete?id=3 -Method PATCH -Headers @{"Authorization"="Bearer $token"; "client_id"="admin" ;"client_secret"="password";"Content-Type"="application/json"}


DELETE TASKS
Invoke-WebRequest -Uri http://localhost:8080/tasks/delete?id=1 -Method DELETE -Headers @{"Authorization"="Bearer $token";"client_id"="admin";"client_secret"="password";"Content-Type"="application/json"}

