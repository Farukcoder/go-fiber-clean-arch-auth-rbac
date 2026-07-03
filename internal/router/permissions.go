package router

// PermissionDefinition describes a single route-level access rule seeded into the database.
type PermissionDefinition struct {
	Name        string
	Method      string
	Path        string
	Description string
}

// PermissionDefinitions returns the full list of route permissions used by the RBAC seeder.
func PermissionDefinitions() []PermissionDefinition {
	return []PermissionDefinition{
		{Name: "auth:login", Method: "POST", Path: "/api/v1/auth/login", Description: "Login user"},
		{Name: "auth:register", Method: "POST", Path: "/api/v1/auth/register", Description: "Register user"},
		{Name: "auth:refresh", Method: "POST", Path: "/api/v1/auth/refresh", Description: "Refresh access token"},
		{Name: "auth:logout", Method: "POST", Path: "/api/v1/auth/logout", Description: "Logout user"},
		{Name: "user:me", Method: "GET", Path: "/api/v1/me", Description: "Get current user"},
		{Name: "log:list", Method: "GET", Path: "/api/v1/logs", Description: "List request logs"},
		{Name: "role:list", Method: "GET", Path: "/api/v1/roles", Description: "List roles"},
		{Name: "role:create", Method: "POST", Path: "/api/v1/roles", Description: "Create role"},
		{Name: "role:update", Method: "PUT", Path: "/api/v1/roles/:id", Description: "Update role"},
		{Name: "role:delete", Method: "DELETE", Path: "/api/v1/roles/:id", Description: "Delete role"},
		{Name: "permission:list", Method: "GET", Path: "/api/v1/permissions", Description: "List permissions"},
		{Name: "permission:create", Method: "POST", Path: "/api/v1/permissions", Description: "Create permission"},
		{Name: "permission:update", Method: "PUT", Path: "/api/v1/permissions/:id", Description: "Update permission"},
		{Name: "permission:delete", Method: "DELETE", Path: "/api/v1/permissions/:id", Description: "Delete permission"},
		{Name: "role_permission:assign", Method: "POST", Path: "/api/v1/roles/:id/permissions", Description: "Assign permission to role"},
		{Name: "role_permission:revoke", Method: "DELETE", Path: "/api/v1/roles/:id/permissions/:permission_id", Description: "Revoke permission from role"},
		{Name: "user:assign_role", Method: "PATCH", Path: "/api/v1/users/:id/role", Description: "Assign role to user"},
	}
}
