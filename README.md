# 📱 Social Simple API

A simple and extensible **social API** project written in Go (with some TypeScript and web frontend parts), designed to provide foundations for building a social backend with REST endpoints and modern tooling. This repository includes server logic, web assets, and DevOps automation for local development.  

---

## 🚀 Overview

**Social Simple API** is an example backend API that implements basic social interactions and services — such as posting content — and demonstrates a clean project structure using Go alongside CLI tools, Docker, and utilities for rapid development.

---

## 📁 Project Structure

```
social-simple-api/
├── .github/ # GitHub workflows
├── cmd/ # Main command / server bootstrap
├── docs/ # Documentation files
├── internal/ # Core business logic and infrastructure
├── web/ # Web frontend and UI assets
├── .air.toml # Air live reload config
├── .envrc.example # Environment example file
├── Makefile # Commands and tasks
├── CHANGELOG.md # Project changes & releases
├── docker-compose.yaml # Docker services for local
├── Dockerfile # Container build definition
├── init.sh # Helper setup script
├── go.mod # Go modules
└── post.http # HTTP API testing snippets
```

---

## 🧠 Tech Stack

This project uses:

- **Go** — for the main backend API  
- **Docker & Docker Compose** — for containerized environment  
- **Make** — to automate tasks  
- **HTTP request snippets** — for quick API testing  
- Web assets and frontend (TypeScript / JS + CSS)

---

## ⚙️ Setup & Local Development

### 1️⃣ Clone the repo

git clone https://github.com/alejandro-cardenas-g/social-simple-api.git
cd social-simple-api
2️⃣ Prepare Environment

Copy the environment example:

cp .envrc.example .env

Edit .env and configure your local settings.

3️⃣ Start Services with Docker
docker compose up -d

This will start any services defined in docker-compose.yaml, such as the database.

4️⃣ Build & Run
go run ./cmd/...

Or use hot reloading if configured (e.g., with Air).

🧪 Testing API Endpoints

Use the post.http file or your preferred REST client (Insomnia / Postman) to quickly test routes such as:

POST /api/posts

GET /api/posts

Auth / login endpoints

Any custom endpoints you add

📄 CHANGELOG

Changes and version history are maintained in CHANGELOG.md.

📈 Possible Improvements

✔ Add authentication
✔ Add user profiles
✔ Add likes, comments, and social interactions
✔ Swagger/OpenAPI documentation
✔ Unit & integration tests
✔ CI/CD pipeline
