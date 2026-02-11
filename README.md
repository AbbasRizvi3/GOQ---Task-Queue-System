# GOQ — Go-based Task Queue System

GOQ is a **task queue and scheduling system** built in Go, designed to support delayed execution, retries with backoff, persistence via PostgreSQL, and a human-friendly web dashboard for observability and control.

---

## 🚀 Features

### Core

* Task scheduling with `run_at` (future execution)
* Persistent task storage using **PostgreSQL**
* Worker-based task processing
* Automatic retries with retry count tracking
* Task state management (`pending`, `ready`, `running`, `failed`, `completed`)

### Backend

* REST APIs for task creation and management
* SQL-backed store (PostgreSQL)
* Authentication with **JWT + cookies**
* Login & signup support
* Dockerized setup for local development

### Web Dashboard (Phase 3)

* Server-Side Rendered (SSR) pages
* Dashboard showing live system state
* Task list and task detail pages
* Live polling via JavaScript
* Retry and cancel task controls
* Task submission via UI
* Safe HTML rendering

### DevOps / CI

* GitHub Actions CI workflow
* Automated builds and tests
* Structured project phases with incremental PRs

---

## 🧱 Architecture Overview

```
Client (Browser / API)
        |
        v
     Go Backend (Gin)
        |
        v
  PostgreSQL (Docker)
        |
        v
   Workers (Goroutines)
```

* **Scheduler** promotes tasks when `run_at <= now`
* **Workers** consume ready tasks
* **SQL store** ensures durability
* **Dashboard** provides real-time visibility and control

---

## 🛠 Tech Stack

* **Language:** Go
* **Web Framework:** Gin
* **Database:** PostgreSQL
* **Auth:** JWT + HTTP Cookies
* **Frontend:** SSR (Go HTML templates) + Vanilla JS
* **Containerization:** Docker & Docker Compose
* **CI:** GitHub Actions

---

## 📦 Project Structure

```
GOQ/
├── cmd/                # Application entry point
├── internal/
│   ├── auth/           # Auth logic (JWT, cookies)
│   ├── handlers/       # HTTP handlers
│   ├── scheduler/      # Task scheduling logic
│   ├── worker/         # Task workers
│   ├── store/          # SQL store (PostgreSQL)
│   └── models/         # Domain models
├── migrations/         # SQL migrations
├── templates/          # SSR HTML templates
├── static/             # JS/CSS assets
├── docker-compose.yml
└── README.md
```

---

## 🐳 Running Locally (Docker)

### Prerequisites

* Docker
* Docker Compose

### Start PostgreSQL

```bash
sudo docker compose up
```

### Stop & Reset Database

```bash
sudo docker compose down
sudo rm -rf ./postgres_data
sudo docker compose up
```

---

## ▶️ Running the Backend

```bash
go run cmd/main.go
```

Server starts on:

```
http://localhost:8000
```

---

## 🌐 Web Dashboard Routes

| Route              | Description     |
| ------------       | --------------- |
| `/api/dashboard`   | Dashboard       |
| `/login`           | Login page      |
| `/signup`          | Signup page     |
| `/api/tasks`       | Task list       |
| `/api/tasks/:id`   | Task detail     |

---

## 📡 API Example

### Create a Task

```http
POST /tasks
Content-Type: application/json
```

```json
{
  "name": "example-task",
  "payload": "hello world",
  "run_at": "2026-02-05T14:30:00Z"
}
```

> `run_at` must be in **RFC3339 / ISO-8601 UTC format**

---

## ⏰ Scheduling Format

GOQ uses **UTC timestamps**:

```
YYYY-MM-DDTHH:MM:SSZ
```

Example:

```
2026-02-05T14:30:00Z
```

---

## 🔁 Retry Behavior

* Tasks automatically retry on failure
* Retry count is tracked per task
* Admins can manually retry or cancel tasks from the dashboard

---

## 🧪 Testing & CI

* Unit tests included for core components
* GitHub Actions runs tests on every PR
* CI ensures build and test stability before merge

---

## 📌 Project Phases

* **Phase 1:** In-memory queue & scheduler
* **Phase 2:** SQL store, auth, Docker, CI
* **Phase 3:** Web dashboard & UX (SSR)
---
