package domain

// RolePermission is the join-table entity linking roles to permissions.
type RolePermission struct {
	RoleID       int64 `json:"role_id" gorm:"column:role_id;primaryKey"`
	PermissionID int64 `json:"permission_id" gorm:"column:permission_id;primaryKey"`
}
