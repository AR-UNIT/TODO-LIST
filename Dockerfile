# Use an official Golang image as a base
FROM golang:1.23 as builder

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache module dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o todo-list-app .

# Start a new image for the runtime (use golang for runtime)
FROM golang:1.23

# Set up a directory for the application
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/todo-list-app .

# Expose the port your app runs on
EXPOSE 8080

# Set the entrypoint to run the application
CMD ["./todo-list-app"]
