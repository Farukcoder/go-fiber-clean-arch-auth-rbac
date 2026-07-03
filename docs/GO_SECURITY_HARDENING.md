# Go Security Hardening Guide — Senior-Level Implementation

This document maps security practices to your existing project (Fiber API
with JWT auth + RBAC) and gives concrete, implementable patterns rather
than abstract advice. Treat this as a checklist to work through, not a
one-time read.

## 1. Authentication

### Password storage

- Use `bcrypt` (cost 12+) or `argon2id` — never MD5/SHA256 directly on
  passwords, and never roll your own hashing.
- Never log a password, even a hashed one, even in debug mode.

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
    return string(hash), err
}

func VerifyPassword(hash, pw string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
```

### JWT — the part most teams get wrong

- **Access tokens: short-lived (5–15 min).** They should be too
  short-lived to matter much if leaked.
- **Refresh tokens: long-lived but rotated on every use, stored
  server-side (hashed) so they can be revoked.** A refresh token that
  can't be revoked is a permanent backdoor if it leaks.
- Sign with a strong algorithm — `RS256`/`ES256` (asymmetric) if the
  token is ever verified by another service; `HS256` is fine for a
  single-service setup but the secret must be high-entropy (32+ random
  bytes) and never hardcoded.
- Always validate `alg` explicitly server-side — never trust the `alg`
  header from the token itself (this is the classic "alg: none" attack).
- Put minimal claims in the token (`user_id`, `role`, `exp`, `iat`) —
  don't put permissions lists or PII in the JWT payload; it's
  base64, not encrypted, and readable by anyone with the token.
- Store refresh tokens **hashed** in your `RefreshToken` table (you
  already have this model) — if the DB leaks, the tokens aren't directly
  usable.

```go
// Refresh token rotation: on every refresh, invalidate the old one
// and issue a new one. This limits the blast radius of a stolen
// refresh token to a single use.
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (*TokenPair, error) {
    hashed := sha256Hex(oldToken)
    stored, err := s.repo.FindRefreshToken(ctx, hashed)
    if err != nil || stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
        return nil, ErrInvalidToken
    }
    // revoke immediately — reuse of an old token is a sign of theft
    if err := s.repo.RevokeToken(ctx, stored.ID); err != nil {
        return nil, err
    }
    return s.issueNewTokenPair(ctx, stored.UserID)
}
```

- If a refresh token is presented twice (reuse after rotation), treat it
  as a signal of compromise and revoke **all** sessions for that user.

### Brute-force / credential stuffing protection

- Rate-limit `/login` per-IP and per-account (see section 8).
- Add a small artificial delay or exponential lockout after repeated
  failures on the same account — but don't lock out indefinitely (that's
  a DoS vector against a legitimate user).
- Consider CAPTCHA after N failed attempts for public-facing login.

## 2. Authorization (RBAC)

You already have Role/Permission/RolePermission tables — the
implementation details that matter:

- **Check permissions server-side on every request, in the handler or
  middleware — never trust a role/permission claim the client can
  influence beyond what was issued in the token.**
- Fetch fresh role/permission data on sensitive actions rather than
  relying solely on what was baked into the JWT at login time — a role
  change should take effect without waiting for token expiry on
  high-privilege actions (admin panel, financial operations).
- Fail closed: if a permission check errors or is ambiguous, deny by
  default.

```go
func RequirePermission(perm string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(string)
        allowed, err := rbacService.HasPermission(c.Context(), userID, perm)
        if err != nil || !allowed {
            return fiber.NewError(fiber.StatusForbidden, "forbidden")
        }
        return c.Next()
    }
}
```

- Apply the principle of least privilege to the seed data too — don't
  seed a default "admin has all permissions" wildcard without an
  explicit, auditable reason.

## 3. Input validation

- Validate every external input at the handler boundary, before it
  reaches `service`. Use `go-playground/validator` with struct tags on
  your `dto` structs — you already have `LoginRequest`,
  `RegisterRequest`, this is where validation tags go.

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=12"`
}
```

- Validate on the server even if the frontend validates too — client
  validation is a UX feature, not a security control.
- Reject unexpected fields where it matters (avoid mass-assignment bugs
  — don't blindly `BodyParser` into a struct that includes fields like
  `Role` or `IsAdmin` that a user shouldn't be able to set).

## 4. SQL injection

- Always use parameterized queries. If you're using `database/sql`
  directly or a query builder, confirm every query uses placeholders,
  never string concatenation or `fmt.Sprintf` with user input.

```go
// Correct
db.QueryContext(ctx, "SELECT * FROM users WHERE email = $1", email)

// Never do this
db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))
```

- If using an ORM (GORM, sqlx, ent), verify raw query helpers
  (`.Raw()`, `.Exec()`) are only ever called with placeholders.

## 5. Transport security & headers

- **Enforce HTTPS/TLS in production** — redirect HTTP to HTTPS, and set
  `Strict-Transport-Security` once TLS is confirmed working.
- Set security headers via middleware (Fiber has `helmet`-equivalent
  patterns):

```go
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Content-Type-Options", "nosniff")
    c.Set("X-Frame-Options", "DENY")
    c.Set("Content-Security-Policy", "default-src 'self'")
    c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
    if !isLocalDev {
        c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
    }
    return c.Next()
})
```

- Lock down CORS explicitly — never `AllowOrigins: "*"` on an
  authenticated API. List actual allowed origins.

```go
app.Use(cors.New(cors.Config{
    AllowOrigins:     "https://yourfrontend.com",
    AllowCredentials: true,
    AllowMethods:     "GET,POST,PUT,DELETE",
}))
```

## 6. Secrets management

- No secrets in source code or committed `.env` files — you already
  have `.env.example` committed and `.env` presumably gitignored;
  double-check `.env` is actually in `.gitignore`, not just assumed.
- In production, load secrets from environment variables injected by
  your orchestrator (Docker secrets, Kubernetes secrets, or a vault
  like HashiCorp Vault / AWS Secrets Manager) — not from a `.env` file
  sitting on the server.
- Rotate the JWT signing secret periodically; support graceful rotation
  by accepting both the old and new secret for a transition window if
  you have long-lived tokens.
- Add a pre-commit hook (`gitleaks` or `trufflehog`) so a secret
  accidentally staged never reaches the remote repo.

## 7. Rate limiting

```go
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP() // or userID for authenticated routes
    },
}))
```

- Apply stricter limits specifically on `/login`, `/register`,
  `/refresh` — these are the endpoints attackers target for brute force
  and enumeration.

## 8. Error handling — don't leak internals

- Never return raw error messages or stack traces to the client in
  production. Map internal errors to generic client-facing messages,
  log the detailed error server-side.

```go
func ErrorHandler(c *fiber.Ctx, err error) error {
    logger.Error("request failed", "path", c.Path(), "err", err)
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": "internal server error", // never err.Error() directly
    })
}
```

- Your `jwt_error_handler.go` should return a generic "invalid or
  expired token" message — don't distinguish "token expired" vs. "token
  malformed" vs. "signature invalid" in the response; that distinction
  helps an attacker probe your token format.

## 9. Logging & audit trail

- Never log passwords, tokens, or full request/response bodies that
  might contain PII or secrets.
- Log security-relevant events with enough context to investigate later:
  failed logins, permission denials, role changes, token revocations —
  with user ID, IP, timestamp, but not the credential itself.
- Your `logger.go` middleware — make sure it redacts `Authorization`
  headers and password fields before logging request data.

```go
func RedactSensitive(headers map[string]string) map[string]string {
    redacted := make(map[string]string, len(headers))
    for k, v := range headers {
        if strings.EqualFold(k, "Authorization") || strings.EqualFold(k, "Cookie") {
            redacted[k] = "[REDACTED]"
            continue
        }
        redacted[k] = v
    }
    return redacted
}
```

## 10. Dependency & code scanning — automate this

Add these to CI so vulnerabilities are caught before merge, not in
production:

```bash
# Known-vulnerability scanner for your actual dependency graph
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Static security analysis (SQL injection patterns, weak crypto, etc.)
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

Add both to `docs/coding/10-code-review-checklist.md` and to your CI
pipeline as required checks, not optional ones.

## 11. Database & infrastructure

- The database user your app connects with should have only the
  privileges it needs (no `DROP`, no access to other schemas) — not the
  Postgres superuser.
- Encrypt DB connections (`sslmode=require` or higher) even on internal
  networks.
- Run the API process as a non-root user in your Docker image; use a
  minimal base image (`distroless` or `alpine`) so there's less attack
  surface if the container is compromised.
- Don't bake `.env` or secrets into the Docker image layer — inject at
  runtime.

## 12. File uploads (if applicable)

- Validate file type by content (magic bytes), not just the extension
  or client-supplied `Content-Type`.
- Enforce a max file size at the framework level before the body is
  fully read into memory.
- Store uploads outside the web root, or in object storage (S3-like)
  with no execute permissions, and serve them through your own handler
  rather than directly from a public directory.

## Priority checklist for your project specifically

Given what's in your current structure, tackle in this order:

1. Confirm refresh tokens are hashed at rest and rotated on use.
2. Confirm `.env` is gitignored and `backend.exe` is untracked (from the
   earlier restructuring pass).
3. Add `govulncheck` and `gosec` to CI.
4. Add rate limiting on `/login`, `/register`, `/refresh`.
5. Add security headers + explicit CORS origin list.
6. Audit `logger.go` and `jwt_error_handler.go` for anything that logs
   or returns sensitive data.
7. Add a permission-check middleware test suite — this is the part most
   likely to have a silent bypass bug.
