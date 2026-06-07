# SkillofIDE Dashboard Backend

This repository contains the microservices backend for the SkillofIDE platform. It consists of multiple Go-based services communicating via gRPC with a custom JSON codec.

---

## 1. Architecture Overview
- **`api-gateway`**: Stateless entrypoint that handles REST requests, parses/validates JWT tokens, and proxies queries.
- **`user-service`**: Microservice that owns the user database and authentication checks.
- **`problem-service`**: Manages practice sets and coding tasks.
- **`submission-service`**: Manages user code submissions.
- **`progress-service`**: Tracks student progress.
- **`notification-service`**: Handles real-time WebSockets.
- **`execution-service`**: Executes code submissions in a sandbox.

---

## 2. Local Setup & Running

To start the database and all microservices, run:
```bash
docker compose up --build -d
```

### Build Required Code Runner Sandbox Images
For the compilation and execution sandbox to work, you must build the docker runner images for each supported language:
```bash
docker build -t skillofide/runner-python:latest services/execution-service/runners/python
docker build -t skillofide/runner-javascript:latest services/execution-service/runners/javascript
docker build -t skillofide/runner-java:latest services/execution-service/runners/java
docker build -t skillofide/runner-cpp:latest services/execution-service/runners/cpp
```

### Port Conflict Troubleshooting (e.g. Metro Bundler)
The React Native Metro Bundler (used for mobile development) runs on port **`8081`** by default. 
To avoid conflicts, the `notification-service` is mapped to host port **`8085`** (binding internally to `8081` in the container). If you need to change other ports, edit [docker-compose.yml](file:///Users/inforbell/skillofied/dashbord_backend/docker-compose.yml).

---

## 3. Database Connection Details
The PostgreSQL database runs inside the Docker container (`dashbord_backend-postgres-1`):
- **Host**: `localhost`
- **Port**: `5432`
- **Database Name**: `skillofide`
- **User**: `skillofide`
- **Password**: `password`

---

## 4. Database Verification

To check if the database exists and view all registered users:

### View Users (without password)
```bash
docker exec -it dashbord_backend-postgres-1 psql -U skillofide -d skillofide -c "SELECT id, email, name, role FROM users;"
```

### View Users (with password)
```bash
docker exec -it dashbord_backend-postgres-1 psql -U skillofide -d skillofide -c "SELECT id, email, name, password, role FROM users;"
```

---

## 5. Adding or Modifying Users (gRPC)

Instead of manually writing SQL insert queries, use the built-in Go utility. It communicates directly with the running `user-service` over gRPC to insert/update the user securely.

```bash
go run add-user.go <email> "<name>" <password> <role>
```

### Examples:
- **Add a student**:
  ```bash
  go run add-user.go student@skillofied.com "Jane Doe" studentpass student
  ```
- **Add/Update an admin**:
  ```bash
  go run add-user.go admin@skillofied.com "Admin User" skillofied123 admin
  ```

---

## 6. Testing the Login API

To verify that the setup is fully functional, make a REST API call to the gateway to retrieve the JWT authentication token:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@skillofied.com","password":"skillofied123"}'
```

**Response Output (JWT token)**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpX...",
  "user": {
    "id": "9f5eff97-b742-4b81-8c0b-0721ba683018",
    "email": "admin@skillofied.com",
    "name": "Admin User",
    "role": "admin"
  }
}
```
