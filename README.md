# Go Fiber Clean Architecture Auth & RBAC Backend

A production-ready Go backend built with Fiber, PostgreSQL, JWT authentication, RBAC permissions, request logging, and security middleware. The project follows a clean architecture split across handlers, services, repositories, routers, and database migrations.

## Features

- JWT authentication with access and refresh tokens
- Role-based access control with permissions and role assignment
- PostgreSQL persistence with migrations and seeders
- Request logging stored in the database
- Rate limiting and security headers middleware
- Static uploads serving from the local storage folder

## Tech Stack

- Go 1.25.11
- Fiber v2
- PostgreSQL + pgx
- GORM
- golang-migrate
- JWT + bcrypt

## Project Structure

```text
cmd/
  api/           # HTTP server entrypoint
  migrate/       # Database migration CLI
  seed/          # Seeder CLI
  debug/         # General debugging utilities
  debug-login/   # Login flow debugging

internal/
  config/        # Environment configuration and logger setup
  domain/        # Core domain models
  dto/           # Request/response DTOs and validation helpers
  handler/       # HTTP handlers
  middleware/    # Auth, logging, RBAC, and security middleware
  repository/    # Database access layer
  router/        # Route registration and permission definitions
  service/       # Business logic

database/
  migrations/    # SQL migration files
  seed/          # Seeders for roles, permissions, and users
```

## Prerequisites

- Go 1.25.11 or newer
- PostgreSQL 12+ (or a compatible server)
- Git
- A `.env` file based on `.env.example`

## Quick Start

1. Clone the repository

```bash
git clone <repository-url>
cd go-fiber-clean-arch-auth-rbac
```

2. Install dependencies

```bash
go mod tidy
```

3. Copy the example environment file

```bash
cp .env.example .env
```

Update the values in `.env` for your local PostgreSQL setup, especially:

```env
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=your_database
DB_USERNAME=your_user
DB_PASSWORD=your_password

JWT_SECRET=replace-with-a-strong-secret
JWT_REFRESH_SECRET=replace-with-a-strong-refresh-secret
```

4. Run database migrations

```bash
go run ./cmd/migrate up
```

5. Seed demo users, roles, and permissions

```bash
go run ./cmd/seed/main.go
```

6. Start the API server

```bash
go run ./cmd/api/main.go
```

The server will start on `http://localhost:8080`.

## Environment Variables

The project uses the values from `.env.example` as defaults. The most important ones are:

```env
APP_NAME=Go App
APP_ENV=local
PORT=8080

DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=database_name
DB_USERNAME=root
DB_PASSWORD=

JWT_SECRET=replace-with-a-strong-secret
JWT_REFRESH_SECRET=replace-with-a-strong-secret-for-refresh-tokens

ALLOWED_ORIGINS=http://localhost:5173,http://localhost:5174,http://localhost:3000
```

## Database Commands

```bash
# Apply all pending migrations
go run ./cmd/migrate up

# Roll back the latest migration
go run ./cmd/migrate down

# Drop and recreate everything from scratch
go run ./cmd/migrate fresh

# Check the current migration version
go run ./cmd/migrate version

# Force the migration version
go run ./cmd/migrate force <version>
```

## Seeders

```bash
# Insert demo users and RBAC configuration
go run ./cmd/seed/main.go
```

## API Overview

### Authentication

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/auth/register` | Create a new user |
| POST | `/api/v1/auth/login` | Authenticate with email and password |
| POST | `/api/v1/auth/refresh` | Rotate access and refresh tokens |
| POST | `/api/v1/auth/logout` | Revoke the current refresh token |
| GET | `/api/v1/me` | Get the authenticated user |

### RBAC and Admin Routes

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/roles` | List all roles |
| POST | `/api/v1/roles` | Create a role |
| PUT | `/api/v1/roles/:id` | Update a role |
| DELETE | `/api/v1/roles/:id` | Delete a role |
| GET | `/api/v1/permissions` | List all permissions |
| GET | `/api/v1/logs` | Retrieve request logs |
| PATCH | `/api/v1/users/:id/role` | Assign a role to a user |

## Default Seeded Credentials

After running the seed command, the following users are created:

| Role | Email | Password |
| --- | --- | --- |
| Super Admin | `superadmin@gmail.com` | `Password123!` |
| Admin | `admin@example.com` | `Password123!` |
| Customer | `demo@example.com` | `DemoPassword1` |

## Development Commands

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Run a specific test package
go test ./internal/service

# Build the API binary
go build -o server ./cmd/api/main.go
```

## Postman

Import the provided `postman_collection.json` file into Postman to explore the API with preconfigured requests.

### Using cURL

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Password123!"}'

# Get current user (replace TOKEN with actual token)
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer TOKEN"
```

## 🐛 Troubleshooting

### Database Connection Error

- Ensure PostgreSQL is running
- Check database credentials in `.env`
- Verify database exists and is accessible

### JWT Token Errors

- Verify `JWT_SECRET` is set correctly
- Check token expiration time
- Ensure `Authorization: Bearer <token>` header format

### Migration Issues

- Check migration files exist in `database/migrations/`
- Verify database has migration table created
- Run `go run ./cmd/migrate/main.go version` to check current version

## 📚 Documentation

- [GO Security Hardening](docs/GO_SECURITY_HARDENING.md)
- [Restructuring Guide](docs/RESTRUCTURE_GUIDE.md)
- [Code Review Checklist](docs/coding/10-code-review-checklist.md)
- [Best Practices](docs/coding/BestPractices.md)

## 📝 License

[Add your license information here]

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
