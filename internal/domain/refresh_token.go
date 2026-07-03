package domain

import (
	"database/sql"
	"time"
)

// RefreshToken is the DB entity for persisted JWT refresh tokens.
type RefreshToken struct {
	ID        int64        `json:"id"`
	UserID    int64        `json:"user_id"`
	JTI       string       `json:"jti"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
