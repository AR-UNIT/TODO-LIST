What is it?

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
In order to provide some guarantees against race conditions for concurrent operations to the db(workers concurrently performing operations on the same rows),
I have written a RowLocking Strategy for queries for write/update operations. 
The below are the implications for row locking strategy with regards to race conditions, and the Read Committed Isolation level set in PostgresDb(default).  

ONLY ROW LOCKING ON WRITES TO DB
THIS WILL NOT BLOCK ANY READS TO ROWS BEEN MODIFIED BY A CONCURRENT PROCESS USING A STANDARD SELECT
THIS ROW LOCKING STRATEGY IMPL IS ALSO USING STANDARD SELECT, ONLY ROW LOCKING FOR MODIFICATIONS
STALE READS COULD OCCUR if a row locked and modified by an update operation, is read using standard select
in Read Committed Isolation level in PostgresSql(default):
/*
 Transaction A locks and updates a row.
 Transaction B reads the row while Transaction A still holds the lock but has not yet committed.
 If Transaction A commits or rolls back, Transaction B may have seen stale or inconsistent data that
	is no longer valid after the commit.
	Problem with ReadCommitted isolation level is non-repeatable reads,
		same row has different values when read at different points of time in the same transaction.
*/

Summary of Race Condition Safety in Read Committed and Row locking for update operations:
Read-Read: Safe, no race condition.
Read-Write: Not safe; may lead to stale reads.
Write-Read: Not safe; may lead to stale or inconsistent reads.
Write-Write: Safe; row-level locking prevents conflicting writes.

 • Architecture Diagram: 


### Distributed TODO API Architecture

1. **Clients**
   - Web or Mobile App communicates with the API over HTTPS.

2. **API Gateway**
   - Handles JWT authentication, rate limiting, and request routing.

3. **Backend Services**
   - Written in Go (or Python for specific components), implements CRUD functionality for TODO items.
   - Adheres to factory and strategy patterns for extensibility.
   - Streams events to Kafka and integrates with Redis for caching.

4. **Kafka**
   - Streams API request events and handles event consumption.

5. **Redis**
   - Acts as a cache to improve query performance.
   - Performs batch updates to PostgreSQL to optimize database writes.

6. **PostgreSQL**
   - Primary database for storing TODO items.

7. **Deployment**
   - Dockerized services orchestrated with Kubernetes.

