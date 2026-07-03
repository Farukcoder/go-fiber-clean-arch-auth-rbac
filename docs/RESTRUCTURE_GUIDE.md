# Migrating This Project to an Idiomatic Go Structure

## Why restructure

The current layout (`app/Http/Controllers`, `app/Http/Middleware`,
`app/Http/Requests`, `app/Models`, `app/Services`) is a direct port of
Laravel's PHP conventions. It works, but it isn't how the Go community
organizes projects, which causes friction for:

- Other Go developers joining the project
- AI coding agents trained on idiomatic Go patterns
- Tooling that assumes Go conventions (linters, `internal/` enforcement)

The fixes below are naming and layering changes only — your actual domain
logic (Auth, RBAC, Logs) doesn't need to change.

## Problems with the current structure

| Current | Issue |
|---|---|
| `app/Http/Controllers`, `Middleware`, `Requests` (PascalCase folders) | Go folder names are always lowercase, no PascalCase/camelCase. This is Laravel naming, not Go. |
| `app/` wrapper directory | Go projects don't use a framework-style root wrapper like this. |
| No `internal/` package | Go's compiler enforces `internal/` — code inside it can't be imported by other modules. Not using it means missing a free architectural guardrail. |
| `app/Models/` and `types/` both hold data structs | Ambiguous — unclear which is a DB entity vs. an API request/response shape. |
| `database/migrate/main.go`, `database/seed/main.go`, `tests/debug/main.go` | Each of these has its own `main.go`, meaning they are separate entrypoints/binaries. They belong under `cmd/`, not under `database/` or `tests/`. |
| `tests/debug`, `tests/debug_login` | These contain `main.go`, not `_test.go` files — they aren't actually tests, they're manual debug scripts. The name is misleading. |
| `backend.exe` tracked in git (shown as modified) | Compiled binaries should never be committed. |
| `routes/`, `repository/`, `config/` live at root while `Services/`, `Models/` live inside `app/` | Inconsistent nesting with no clear rule for what goes where. |

## Target structure

```
backend/
├── cmd/
│   ├── api/main.go              # was: root main.go
│   ├── migrate/main.go          # was: database/migrate/main.go
│   ├── seed/main.go             # was: database/seed/main.go
│   ├── debug/main.go            # was: tests/debug/main.go
│   └── debug-login/main.go      # was: tests/debug_login/main.go
│
├── internal/
│   ├── handler/                 # was: app/Http/Controllers
│   │   ├── auth_handler.go
│   │   ├── logs_handler.go
│   │   └── rbac_handler.go
│   ├── middleware/               # was: app/Http/Middleware
│   │   ├── jwt_error_handler.go
│   │   ├── logger.go
│   │   └── rbac.go
│   ├── dto/                       # was: app/Http/Requests + types/
│   │   ├── login_request.go
│   │   ├── register_request.go
│   │   ├── auth_types.go
│   │   ├── rbac_types.go
│   │   └── response.go
│   ├── domain/                     # was: app/Models
│   │   ├── log.go
│   │   ├── permission.go
│   │   ├── refresh_token.go
│   │   ├── role.go
│   │   ├── role_permission.go
│   │   ├── user.go
│   │   └── user_role.go
│   ├── service/                     # was: app/Services
│   │   ├── auth_service.go
│   │   └── rbac_service.go
│   ├── repository/                   # unchanged, just nested under internal/
│   ├── router/                        # was: routes/
│   └── config/                         # was: config/
│
├── database/
│   ├── migrations/               # *.sql files stay at top level — not code
│   ├── database.go
│   └── seed/seeders/              # rbac_seeder.go, user_seeder.go
│
├── pkg/
│   ├── httpclient/                # was: http_services/ (if it wraps external API calls)
│   └── response/                    # standard JSON response helper
│
├── docs/coding/BestPractices.md
├── storage/logs/                    # keep, but gitignore the log files
├── .env / .env.example
├── AGENTS.md / CLAUDE.md / GEMINI.md
├── postman_collection.json
├── go.mod / go.sum
└── README.md
```

### Two folders that need your judgment

- **`app/Resources/`** — contents weren't visible in the file tree, but by
  Laravel convention this usually holds API response transformers. Likely
  maps to `internal/dto/` or `pkg/response/`. Check what's actually
  inside before moving it.
- **`http_services/`** — the name suggests external API client wrappers,
  so it's mapped to `pkg/httpclient/`. If it actually contains internal
  business logic instead, keep it under `internal/` instead.

## Layer responsibility rules

- `handler/` — parses the request, calls `service`, writes the response.
  No business logic here.
- `service/` — business logic and use cases. Never imports the web
  framework (`fiber.Ctx`, `http.Request`, etc.) directly — this keeps it
  framework-independent and easy to unit test.
- `repository/` — implements interfaces defined by `domain/` or
  `service/`, talks to the database.
- `domain/` — entities and core business rules. Imports nothing else
  from this project.
- `dto/` — request/response shapes for the API boundary, separate from
  `domain/` entities so your DB model and API contract can evolve
  independently.

Dependency direction: `handler → service → repository → domain`, never
the reverse.

## Migration steps

1. **Move files with `git mv`**, not plain `mv` — this preserves file
   history. Group moves by layer (all controllers, then all middleware,
   etc.) so each commit is reviewable.
2. **Fix package declarations.** Every moved file currently has something
   like `package Controllers` or `package Models`. Update these to match
   the new lowercase package name (`package handler`, `package domain`,
   etc.).
3. **Fix import paths** across the project to point at the new package
   locations.
4. **Stop tracking `backend.exe`:**
   ```bash
   git rm --cached backend.exe
   echo "backend.exe" >> .gitignore
   echo "storage/logs/*.log" >> .gitignore
   ```
5. **Rebuild and fix errors:**
   ```bash
   go build ./...
   go vet ./...
   ```
6. **Review before committing:**
   ```bash
   git status --short
   ```

## Doing this with an AI coding agent

Steps 2–5 (package names, import paths, fixing compile errors) are
mechanical and well-suited to an AI agent — point it at the compiler
errors after the file moves and it can resolve them systematically. Do
the `git mv` file moves yourself (or via a script) first, since that's
the part where you want full control over exactly where each file ends
up; then hand the compile-fixing loop to the agent.

## Reference

This target layout follows the community `golang-standards/project-layout`
conventions (not an official Go team standard, but the most widely
adopted pattern in the Go ecosystem), scaled down to what a
medium-sized API project actually needs.
