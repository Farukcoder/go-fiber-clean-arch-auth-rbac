# Code Review Security & Vulnerability Checklist

Prior to merging any pull request or deploying code to staging/production, perform the following verification checks:

## 1. Local Security Scans
Developers must run these scans locally:

```bash
# 1. Check for known vulnerabilities in your actual dependency graph
govulncheck ./...

# 2. Run static security analysis (SQL injection, weak cryptography, etc.)
gosec ./...
```

Ensure both commands exit successfully with zero issues. If any issue is reported, fix it immediately or document the explicit security exception with lead approval.

## 2. Security Design Checklist
Verify that:
- **Authentication:** No credentials, hashes, or sensitive tokens are printed in logging libraries or saved in plaintext databases.
- **Refresh Token Rotation:** Refresh tokens are hashed using SHA-256 before database insertion, and compromised/reused tokens trigger a complete session revocation.
- **Authorization:** Permissions are checked server-side and default to forbidden (fail-closed).
- **Input Validation:** Every payload schema has strict struct tags (e.g. `validate:"required,email"`).
- **SQL Injection:** Avoid `fmt.Sprintf` or string concatenation in database queries; use GORM's parameterized query engine.
- **Security Headers:** All HTTP routes serve standard headers (`X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Content-Security-Policy`, `Referrer-Policy`).
- **CORS Configuration:** Explicit origin lists must be configured in `.env` (`ALLOWED_ORIGINS`). `AllowOrigins: "*"` is prohibited in staging or production.
