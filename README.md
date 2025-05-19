# FX gRPC Task List Application

This project is a Task List application built with Go, featuring a gRPC service for task management. It utilizes Uber FX for dependency injection, Cobra for CLI commands, and `database/sql` with a MySQL driver for database interaction. This README provides a comprehensive guide to understanding, building, running, and testing the application.

## Table of Contents

- [Features](#features)
- [Key Packages Used](#key-packages-used)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Code Generation](#code-generation)
- [Build Instructions](#build-instructions)
  - [Local Build](#local-build)
  - [Docker Build](#docker-build)
- [Running the Application](#running-the-application)
  - [Database Setup (MySQL)](#database-setup-mysql)
  - [Running the Server (Local)](#running-the-server-local)
  - [Running the Server (Docker)](#running-the-server-docker)
- [Using the Client CLI](#using-the-client-cli)
- [Interacting with the API](#interacting-with-the-api)
  - [gRPC API](#grpc-api)
- [Error Handling and Logging](#error-handling-and-logging)

## Features

- gRPC service (`TaskService`) for managing tasks:
  - `AddTask(title, description, status)`: Adds a new task.
  - `GetTasks()`: Retrieves a list of all tasks.
  - `CompleteTask(task_id)`: Marks an existing task as completed.
- CLI client to interact with the gRPC service's functionalities.
- Dependency injection managed by Uber FX.
- MySQL database interaction using the standard `database/sql` package.
- Configuration primarily through environment variables with sensible defaults.
- Docker support for easy containerization and deployment.

## Key Packages Used

### Standard Library

- **`context`**: For managing deadlines, cancellation signals, and request-scoped values.
- **`database/sql`**: For generic SQL (or SQL-like) database interaction.
- **`fmt`**: For formatted I/O operations.
- **`os`**: For operating system functionalities like reading environment variables and signal handling.
- **`time`**: For time-related operations and timeouts.

### Third-Party Libraries

- **`github.com/spf13/cobra`**: A powerful library for creating modern CLI applications.
- **`go.uber.org/fx`**: A dependency injection framework from Uber.
- **`go.uber.org/zap`**: A structured logging library from Uber.
- **`google.golang.org/grpc`**: The official Go implementation of gRPC.
- **`github.com/go-sql-driver/mysql`**: The MySQL driver for Go's `database/sql` package.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go:** Version 1.21 or later.
- **Docker & Docker Compose:** For containerized setup and deployment.
- **MySQL Server:** A running instance (local or Dockerized) for the database.
- **`protoc` Compiler:** The Protocol Buffers compiler.
- **Go gRPC Plugins:**
  - `protoc-gen-go`: For generating Go protobuf structs.
  - `protoc-gen-go-grpc`: For generating Go gRPC client and server stubs.

You can install the Go gRPC plugins using:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Project Structure

```
.
├── api/                     # Generated protobuf files
│   ├── api.pb.go
│   ├── api_grpc.pb.go
├── cmd/                     # CLI commands
│   ├── addTask.go
│   ├── client.go
│   ├── completeTask.go
│   ├── getTasks.go
│   ├── root.go
│   ├── server.go
├── client/                  # gRPC client setup
│   └── client.go
├── config/                  # Configuration management
│   └── config.go
├── database/                # Database connection setup
│   └── database.go
├── docker/                  # Docker-related files
│   ├── Dockerfile
│   ├── docker-compose.yml
├── init-db/                 # Database initialization script
│   └── init-db.sh
├── repository/              # Task repository for database operations
│   └── task_repository.go
├── server/                  # gRPC server and service implementation
│   ├── api_service.go
│   ├── server.go
├── main.go                  # Entry point for the application
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies checksum
├── README.md                # Project documentation
└── .gitignore               # Git ignore file
```

## Configuration

The application is configured using environment variables. Sensible defaults are provided:

- `DB_USER`: MySQL username (default: `user`)
- `DB_PASSWORD`: MySQL password (default: `password`)
- `DB_HOST`: MySQL host (default: `localhost`)
- `DB_PORT`: MySQL port (default: `3306`)
- `DB_NAME`: MySQL database name (default: `taskdb`)
- `GRPC_PORT`: Port for the gRPC server (default: `50051`)

## Code Generation

To generate Go code from your `.proto` files, run:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api.proto
```

## Build Instructions

### Local Build

To build the application binary locally:

```bash
go build -o fx-grpc-app .
```

### Docker Build

To build the Docker image:

```bash
docker build -t fx-grpc-tasklist-app -f docker/Dockerfile .
```

## Running the Application

### Database Setup (MySQL)

Ensure you have a MySQL server running. You can use the provided `docker-compose.yml` to start one:

```bash
docker-compose -f docker/docker-compose.yml up -d mysql-db
```

### Running the Server (Local)

After building the application:

```bash
./fx-grpc-app server
```

### Running the Server (Docker)

Using Docker Compose:

```bash
docker-compose -f docker/docker-compose.yml up grpc-server
```

## Using the Client CLI

The application includes a CLI client built with Cobra to interact with the gRPC service. First, ensure the server is running.

### Add a Task

```bash
./fx-grpc-app client add-task --title "My Task" --description "Task description" --status "pending"
```

### Get All Tasks

```bash
./fx-grpc-app client get-tasks
```

### Complete a Task

```bash
./fx-grpc-app client complete-task --id <task_id>
```

## Interacting with the API

### gRPC API

The gRPC API is defined in `api.proto`. The `TaskService` exposes the following methods:

- `GetTasks(GetTasksRequest) returns (GetTasksReply)`
- `AddTask(AddTaskRequest) returns (AddTaskReply)`
- `CompleteTask(CompleteTaskRequest) returns (CompleteTaskReply)`

## Error Handling and Logging

The application uses `zap` for structured logging. Logs are output to standard output. gRPC errors are returned with appropriate gRPC status codes.