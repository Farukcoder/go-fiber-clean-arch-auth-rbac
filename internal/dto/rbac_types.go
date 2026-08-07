package dto

// RolePermissionAssignment is used when assigning a permission to a role.
type RolePermissionAssignment struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

// CreateRoleInput is the request body for creating a role.
type CreateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRoleInput is the request body for updating a role.
type UpdateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreatePermissionInput is the request body for creating a permission.
type CreatePermissionInput struct {
	Name        string `json:"name"`
	Module      string `json:"module"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

// UpdatePermissionInput is the request body for updating a permission.
type UpdatePermissionInput struct {
	Name        string `json:"name"`
	Module      string `json:"module"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

// AssignPermissionInput is the request body for assigning a permission to a role.
type AssignPermissionInput struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

// AssignRoleInput is the request body for assigning a role to a user.
type AssignRoleInput struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}
