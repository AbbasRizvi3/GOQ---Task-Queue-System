# GOQ — Asynchronous Task Queue System in Go

## Overview

GOQ is a self-hosted asynchronous task queue system written in Go.  
This project is a simplified assignment implementation focusing on **concurrent task processing, worker pools, retries, and a simple REST API**.

The system demonstrates core Go concepts including:

- Goroutines & channels
- Worker pools
- Task queue management
- Error handling and retries with backoff
- Context-based cancellation
- Simple in-memory task storage
- REST API endpoints

---

## Features To Implement

- **Task Queue (in-memory)**  
  - Add tasks dynamically
  - Tasks stored in memory slice or channel

- **Worker Pool**  
  - Configurable number of workers
  - Concurrent task execution
  - Logging task status

- **Retries with Backoff**  
  - Retry failed tasks with configurable delay
  - Exponential backoff strategy

- **REST API Endpoints**  
  - `POST /tasks` → Submit a new task  
  - `GET /tasks` → List all tasks  
  - `GET /tasks/{id}` → Get task by ID  

- **Logging**  
  - Task execution logs with timestamp, status, and result

- **Testing**  
  - Unit tests for task processing, queue operations, and worker pool

---

## Project Structure

```

goq/
├─ cmd/
│  ├─ goq-server/       # entry point for the API server
├─ internal/
│  ├─ app/              # dependency wiring
│  ├─ queue/            # in-memory task queue
│  ├─ worker/           # worker pool
│  ├─ domain/           # task structs & interfaces
│  ├─ server/           # REST API implementation
├─ tasks.log            # sample task execution logs
├─ go.mod
├─ go.sum
└─ README.md

---

## Getting Started

### Prerequisites

- Go 1.21+ installed
- (Optional) Docker if you want to containerize

### Run Locally

1. Clone the repository:

```bash
git clone <your-repo-url>
cd goq
````

2. Run the server:

```bash
go run cmd/goq-server/main.go
```

3. Submit and list tasks via REST API:

```bash
# Add a task
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"name":"Task 1"}'

# List tasks
curl http://localhost:8080/tasks
```

---

## Testing

Run all unit tests with:

```bash
go test ./... -v -race
```

---
