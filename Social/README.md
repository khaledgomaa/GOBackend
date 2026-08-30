# Social API

A Go backend API built with `chi`, PostgreSQL, and Docker.

## Prerequisites
- [Go](https://go.dev/)
- [Docker](https://www.docker.com/) & Docker Compose
- `make`
- [Air](https://github.com/cosmtrek/air) (for live-reloading)
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

## Getting Started

1. **Start the database:**
   Run a PostgreSQL instance using Docker Compose:
   ```bash
   docker-compose up -d
   ```

2. **Environment Setup:**
   Create a `.env` file in the root directory. The application has default fallbacks, so an empty file is sufficient to satisfy the `Makefile` requirement.

3. **Run Database Migrations:**
   Apply the database schema:
   ```bash
   make migrate-up
   ```

4. **Run the Application:**
   Start the application with live-reloading (via Air):
   ```bash
   make run
   ```
   Or, build the application without running it:
   ```bash
   make build
   ```

## Makefile Reference

The project includes a `Makefile` with commands to simplify common tasks:

- `make run`: Runs the application with live-reloading using `air`.
- `make build`: Compiles the Go application and places the executable at `./bin/main`.
- `make migrate-up`: Applies all pending database migrations.
- `make migrate-down`: Rolls back the last applied database migrations.
- `make migrate-create name=<migration_name>`: Scaffolds a new migration file.

## Testing the API

The project includes an `api.http` file for quickly testing endpoints directly from VS Code.

**To use it:**
1. Install the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension in VS Code.
2. Open the `api.http` file in the root directory.
3. Click the **Send Request** text button that appears above each endpoint to execute it.
