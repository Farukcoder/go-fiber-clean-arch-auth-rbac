package domain

import "time"

// User is the DB entity for application users.
type User struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey"`
	Name         string    `json:"name" gorm:"column:name"`
	Email        string    `json:"email" gorm:"column:email"`
	Phone        string    `json:"phone" gorm:"column:phone"`
	RoleID       int64     `json:"role_id" gorm:"column:role_id"`
	RoleName     string    `json:"role_name" gorm:"column:role_name"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
}
