package router

// PermissionDefinition describes a single route-level access rule seeded into the database.
type PermissionDefinition struct {
	Name        string
	Module      string
	Method      string
	Path        string
	Description string
}

// PermissionDefinitions returns the full list of route permissions used by the RBAC seeder.
func PermissionDefinitions() []PermissionDefinition {
	return []PermissionDefinition{
		{Name: "auth:login", Module: "AUTH", Method: "POST", Path: "/api/v1/auth/login", Description: "Login user"},
		{Name: "auth:register", Module: "AUTH", Method: "POST", Path: "/api/v1/auth/register", Description: "Register user"},
		{Name: "auth:refresh", Module: "AUTH", Method: "POST", Path: "/api/v1/auth/refresh", Description: "Refresh access token"},
		{Name: "auth:logout", Module: "AUTH", Method: "POST", Path: "/api/v1/auth/logout", Description: "Logout user"},
		{Name: "user:me", Module: "USERS", Method: "GET", Path: "/api/v1/me", Description: "Get current user"},
		{Name: "user:me_permissions", Module: "USERS", Method: "GET", Path: "/api/v1/me/permissions", Description: "Get current user's permissions"},
		{Name: "user:list", Module: "USERS", Method: "GET", Path: "/api/v1/users", Description: "List all users"},
		{Name: "log:list", Module: "SYSTEM", Method: "GET", Path: "/api/v1/logs", Description: "List request logs"},
		{Name: "role:list", Module: "RBAC", Method: "GET", Path: "/api/v1/roles", Description: "List roles"},
		{Name: "role:create", Module: "RBAC", Method: "POST", Path: "/api/v1/roles", Description: "Create role"},
		{Name: "role:update", Module: "RBAC", Method: "PUT", Path: "/api/v1/roles/:id", Description: "Update role"},
		{Name: "role:delete", Module: "RBAC", Method: "DELETE", Path: "/api/v1/roles/:id", Description: "Delete role"},
		{Name: "role:permissions:list", Module: "RBAC", Method: "GET", Path: "/api/v1/roles/:id/permissions", Description: "Get role permissions"},
		{Name: "permission:list", Module: "RBAC", Method: "GET", Path: "/api/v1/permissions", Description: "List permissions"},
		{Name: "permission:create", Module: "RBAC", Method: "POST", Path: "/api/v1/permissions", Description: "Create permission"},
		{Name: "permission:update", Module: "RBAC", Method: "PUT", Path: "/api/v1/permissions/:id", Description: "Update permission"},
		{Name: "permission:delete", Module: "RBAC", Method: "DELETE", Path: "/api/v1/permissions/:id", Description: "Delete permission"},
		{Name: "role_permission:assign", Module: "RBAC", Method: "POST", Path: "/api/v1/roles/:id/permissions", Description: "Assign permission to role"},
		{Name: "role_permission:revoke", Module: "RBAC", Method: "DELETE", Path: "/api/v1/roles/:id/permissions/:permission_id", Description: "Revoke permission from role"},
		{Name: "user:assign_role", Module: "USERS", Method: "PATCH", Path: "/api/v1/users/:id/role", Description: "Assign role to user"},
		{Name: "user:permissions:list", Module: "USERS", Method: "GET", Path: "/api/v1/me/permissions", Description: "Get user permissions"},
	}
}
