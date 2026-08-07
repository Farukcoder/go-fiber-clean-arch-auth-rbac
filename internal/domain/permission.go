package domain

import "time"

// Permission is the DB entity for RBAC permissions (route-level access rules).
type Permission struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey"`
	Name        string    `json:"name" gorm:"column:name"`
	Module      string    `json:"module" gorm:"column:module"`
	Method      string    `json:"method" gorm:"column:method"`
	Path        string    `json:"path" gorm:"column:path"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}
