package domain

import "time"

// UserRole is the join-table entity linking users to their assigned role.
type UserRole struct {
	UserID    int64     `json:"user_id" gorm:"column:user_id;primaryKey"`
	RoleID    int64     `json:"role_id" gorm:"column:role_id;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
