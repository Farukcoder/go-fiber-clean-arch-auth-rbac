package dto

import "time"

// UserResponse is the API representation of a user (excludes sensitive fields).
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	RoleID    int64     `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}
