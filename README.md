# Go Fiber Clean Architecture Auth & RBAC Backend

A production-ready Go backend using Fiber framework with clean architecture. Features JWT authentication, role-based access control (RBAC), refresh token rotation, PostgreSQL database, request logging, rate limiting, security hardening, database migrations, and automated seeders.

## ✨ Features

- **JWT Authentication** - Secure token-based authentication with refresh token rotation
- **Role-Based Access Control (RBAC)** - Granular permission management with roles and permissions
- **Clean Architecture** - Organized layered structure for maintainability and scalability
- **PostgreSQL Integration** - Robust database with migration and seeding support
- **Request Logging** - Comprehensive request/response logging with database persistence
- **Rate Limiting** - Built-in rate limiter to prevent abuse
- **Security Middleware** - CORS, CSRF protection, and security hardening
- **Database Migrations** - Version-controlled schema management
- **Automated Seeders** - Demo data and RBAC setup

## 📚 Tech Stack

- **Framework**: Fiber (Express-like Go web framework)
- **Database**: PostgreSQL 12+
- **Authentication**: JWT (JSON Web Tokens)
- **Language**: Go 1.25+
- **Migration Tool**: golang-migrate

## Directory Structure

```
backend/
├── cmd/
│   ├── api/                  # Main HTTP server entrypoint
│   ├── migrate/              # Migration tool entrypoint
│   ├── seed/                 # Seeder tool entrypoint
│   ├── debug/                # Debug script
│   └── debug-login/          # Login flow debug script
│
├── internal/                 # Private application and business logic
│   ├── config/               # App configuration logic
│   ├── domain/               # Domain model structs (RequestLog, User, Role, etc.)
│   ├── dto/                  # Data Transfer Objects & validation
│   ├── handler/              # HTTP request handlers (controllers)
│   ├── middleware/           # HTTP middleware logic
│   ├── repository/           # Data repository layer (database operations)
│   ├── router/               # Route setup & permission definitions
│   └── service/              # Business logic services
│
├── database/                 # SQL migrations and database connection setup
│   ├── migrations/           # SQL migration files
│   └── database.go           # Database connection logic
│
├── pkg/                      # Shared helper packages
│   └── httpclient/           # External HTTP client wrappers
│
├── storage/                  # Local filesystem storage / logs
├── .env                      # Local configuration file (not committed)
├── .env.example              # Example configuration template
├── go.mod                    # Go module dependencies
└── go.sum                    # Go module checksums
```

## 🔌 API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|----------------|
| POST | `/api/v1/auth/login` | User login with email/password | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh access and refresh tokens | ✅ |
| POST | `/api/v1/auth/logout` | Revoke the current refresh token | ✅ |
| GET | `/api/v1/me` | Get current authenticated user | ✅ |

### RBAC Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|----------------|
| GET | `/api/v1/roles` | List all roles | ✅ |
| GET | `/api/v1/permissions` | List all permissions | ✅ |
| POST | `/api/v1/users/:id/roles` | Assign role to user | ✅ |

### Logging Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|----------------|
| GET | `/api/v1/logs` | Retrieve request logs | ✅ |

### Standard Response Format

All API responses follow this consistent structure:

```json
{
  "status": true,
  "status_code": 200,
  "message": "Success message here",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

### Error Response Format

```json
{
  "status": false,
  "status_code": 400,
  "message": "Error description",
  "data": null
}
```

## 🚀 Quick Start

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- Git
- (Optional) Postman for API testing

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-fiber-clean-arch-auth-rbac
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and settings
   ```

4. **Run database migrations**
   ```bash
   go run ./cmd/migrate/main.go up
   ```

5. **Seed database with demo data**
   ```bash
   go run ./cmd/seed/main.go
   ```

6. **Start the development server**
   ```bash
   go run ./cmd/api/main.go
   ```
   Server runs on `http://localhost:8080`

## 🗄️ Database Management

### Migration Commands

```bash
# Apply all pending migrations
go run ./cmd/migrate/main.go up

# Drop all tables and recreate from scratch
go run ./cmd/migrate/main.go fresh

# Roll back the latest migration
go run ./cmd/migrate/main.go down

# Check current migration version
go run ./cmd/migrate/main.go version

# Force set migration version
go run ./cmd/migrate/main.go force <version>
```

### Seeding Commands

```bash
# Insert demo users and RBAC configuration
go run ./cmd/seed/main.go
```

### Debug Commands

```bash
# Run general debug operations
go run ./cmd/debug/main.go

# Debug login flow with test user
go run ./cmd/debug-login/main.go
```

## 🛠️ Development

### Build the Application

```bash
# Build for current OS
go build -o server ./cmd/api/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server ./cmd/api/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o server.exe ./cmd/api/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Check for common mistakes
golangci-lint run
```

## 👥 Default Credentials (After Seeding)

| User Type | Email | Password |
|-----------|-------|----------|
| Super Admin | superadmin@gmail.com | Password123! |
| Admin | admin@example.com | Password123! |
| Demo | demo@example.com | DemoPassword1 |

## 🔐 Security Features

- **JWT Authentication** - Secure token-based authentication with configurable expiration
- **Refresh Token Rotation** - Automatic token refresh with secure rotation
- **Password Hashing** - Bcrypt for secure password storage
- **CORS Configuration** - Configurable cross-origin resource sharing
- **Rate Limiting** - Prevent brute force attacks with rate limiting
- **Request Validation** - Comprehensive input validation and sanitization
- **Security Headers** - HTTP security headers middleware
- **Error Handling** - Secure error messages without leaking sensitive information

## 🛡️ Security Scanning & Vulnerability Management

### Running Vulnerability Checks

The project uses `govulncheck` to scan for known vulnerabilities in dependencies and the Go standard library.

```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run vulnerability scan
gvulncheck ./...

# For verbose output
gvulncheck -show verbose ./...
```

### Keeping Go Updated

Vulnerabilities in the standard library are patched in Go release updates. Ensure you're running the latest stable version:

```bash
# Check your current Go version
go version

# Download latest Go: https://go.dev/dl/
# Update go.mod to the latest Go version
# (current: 1.25.11 or higher)
```

**Note:** This project requires **Go 1.25.11 or higher** to address all standard library vulnerabilities.

### Dependency Vulnerability Management

Regularly update dependencies:

```bash
# Get latest compatible versions
go get -u ./...

# Update specific module
go get -u github.com/module/name

# Clean up
go mod tidy
```

### CI/CD Integration

Add to your CI/CD pipeline to catch vulnerabilities automatically:

```yaml
# GitHub Actions example
- name: Run Vulnerability Check
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
```

## 📁 Project Structure Explanation

```
internal/
├── config/         # Configuration management and logger setup
├── domain/         # Business entities and models
├── dto/            # Data Transfer Objects with validation rules
├── handler/        # HTTP request handlers (controllers)
├── middleware/     # HTTP middleware (logging, RBAC, security)
├── repository/     # Data access layer (database operations)
├── router/         # Route definitions and permission mapping
└── service/        # Business logic layer

cmd/
├── api/            # Production server entrypoint
├── migrate/        # Database migration tool
├── seed/           # Database seeder tool
├── debug/          # Development debug utilities
└── debug-login/    # Login flow debugging

database/
├── migrations/     # SQL migration files with version control
└── seed/           # Seeders for demo data
```

## 🚢 Deployment

### Environment Variables

See `.env.example` for all required environment variables:

```bash
# Core
PORT=8080
ENV=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_rbac

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=900

# Redis (optional)
REDIS_URL=localhost:6379
```

### Building for Production

```bash
go build -ldflags="-w -s" -o server ./cmd/api/main.go
```

### Docker Support (if configured)

```bash
docker build -t go-fiber-auth .
docker run -p 8080:8080 --env-file .env go-fiber-auth
```

## 📊 Database Schema

The project includes migrations for:

- **users** - User accounts with credentials
- **roles** - Role definitions for RBAC
- **permissions** - Fine-grained permission definitions
- **user_roles** - User-to-role assignments
- **role_permissions** - Role-to-permission mappings
- **refresh_tokens** - Token refresh tracking
- **logs** - Request/response logging

## 🧪 Testing the API

### Using Postman

Import `postman_collection.json` into Postman to test all endpoints with pre-configured requests.

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
