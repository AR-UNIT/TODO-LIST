This is a project I used to learn backend implementation and principles in GoLang.
The project intends to implement logically simple CRUD operations of a TODO application, however, tries to include
as many aspects and tools/frameworks as possible in the implementation.

•	PostgresDB as the default storage method, with code open to extension to different storage methods, and an interface defined to implement a task manager

•	Application containerized using Docker & Kubernetes as follows:
 - POSTGRES CONTAINER
 - ZOOKEEPER CONTAINER
 - KAFKA CONTAINER
 - REDIS CONTAINER
 - APP CONTAINER

 • Middlewares implemented:
   - JWT authentication for every api call + autentication endpoint to get JWT tokens
   - rate limiting for every end point
   - Kafka for streaming events from API requests and consuming events at task manager, decoupled api hit with processing the request
   - Redis for caching and performing batch database updates, to reduce number of calls made to DB as the number of application nodes grows


 • DB Isolation level and guarantees against concurrent operations:

 • Architecture Diagram: 
