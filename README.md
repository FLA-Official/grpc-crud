# grpc-crud

A simple gRPC CRUD example in Go with a REST gateway.

## Overview

This project demonstrates a basic CRUD service for `User` records using:

- gRPC server (`cmd/server/main.go`)
- gRPC Gateway REST proxy (`gateway/main.go`)
- PostgreSQL persistence with `sqlx`
- Clean separation of handler, service, and repository layers
- Protocol buffers definitions in `proto/user/v1/user.proto`

## Architecture

- `cmd/server/main.go`: starts the gRPC server on port `50051`
- `gateway/main.go`: starts the HTTP REST gateway on port `8080` and forwards requests to gRPC
- `internal/config/config.go`: loads database configuration from `.env`
- `internal/user/handler/user_handler.go`: implements the gRPC service methods
- `internal/user/service/user_service.go`: business logic and validation
- `internal/user/repo/user_repo.go`: database operations using `sqlx`
- `internal/user/model/user_model.go`: domain model for `User`
- `proto/user/v1/user.proto`: proto service and message definitions
- `pkg/pb/userpb`: generated Go code from protobuf definitions

## Data Model

The `User` model stores:

- `id` (int64)
- `name` (string)
- `email` (string)

## Expected Database Schema

The repository expects a PostgreSQL table named `users`.

Example schema:

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL
);
```

## Environment Variables

The server uses environment variables loaded from `.env` via `github.com/joho/godotenv`.

Required variables:

- `DB_USER`
- `DB_PASSWORD`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`

Example `.env`:

```env
DB_USER=postgres
DB_PASSWORD=secret
DB_HOST=localhost
DB_PORT=5432
DB_NAME=grpc_crud_db
```

## gRPC API

The gRPC service is defined in `proto/user/v1/user.proto` as `UserService`.

Implemented RPC methods:

- `CreateUser(CreateUserRequest) returns (CreateUserResponse)`
- `GetUser(GetUserRequest) returns (GetUserResponse)`
- `UpdateUser(UpdateUserRequest) returns (UpdateUserResponse)`
- `DeleteUser(DeleteUserRequest) returns (DeleteUserResponse)`

There is also a `Greetings(HelloRequest)` RPC defined in the proto, but it is not implemented in `internal/user/handler/user_handler.go`.

## REST Gateway Endpoints

The gRPC gateway exposes HTTP/JSON endpoints mapped from the proto annotations:

- `GET /v1/hello/{name}` -> `Greetings`
- `POST /v1/users` -> `CreateUser`
- `GET /v1/users/{id}` -> `GetUser`
- `PUT /v1/users/{id}` -> `UpdateUser`
- `DELETE /v1/users/{id}` -> `DeleteUser`

## Running the Project

1. Install dependencies:

   ```bash
go mod tidy
   ```

2. Start PostgreSQL and create the `users` table.

3. Populate `.env` with database credentials.

4. Run the gRPC server:

   ```bash
go run cmd/server/main.go
   ```

5. In another terminal, run the REST gateway:

   ```bash
go run gateway/main.go
   ```

6. Access the REST API on `http://localhost:8080`.

## Example REST Requests

Create a user:

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

Get a user:

```bash
curl http://localhost:8080/v1/users/1
```

Update a user:

```bash
curl -X PUT http://localhost:8080/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe Updated","email":"jane.new@example.com"}'
```

Delete a user:

```bash
curl -X DELETE http://localhost:8080/v1/users/1
```

## gRPC Testing

You can also call the gRPC endpoint directly on `localhost:50051` using tools like `grpcurl`.

Example:

```bash
grpcurl -plaintext -d '{"name":"Jane","email":"jane@example.com"}' localhost:50051 user.v1.UserService.CreateUser
```

## Notes

- The project uses `grpc-gateway` to expose a RESTful interface for the gRPC service.
- Database access is handled through `sqlx` and the PostgreSQL driver `lib/pq`.
- The proto-generated Go package path is `pkg/pb/userpb`.
- The `Greetings` method exists in the proto file but is not handled in the service implementation.

## File Overview

- `cmd/server/main.go`: gRPC server bootstrap
- `gateway/main.go`: REST gateway bootstrap
- `internal/config/config.go`: env configuration loader
- `internal/user/handler/user_handler.go`: gRPC service handler
- `internal/user/service/user_service.go`: business logic
- `internal/user/repo/user_repo.go`: database repository
- `internal/user/model/user_model.go`: user entity
- `proto/user/v1/user.proto`: protobuf definitions
- `pkg/pb/userpb/`: generated protobuf Go code
