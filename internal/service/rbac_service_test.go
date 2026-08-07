package service

import (
	"testing"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
)

func TestRBACService_Allowed(t *testing.T) {
	tests := []struct {
		name       string
		roleID     int64
		method     string
		path       string
		snapshot   rolePermissionsSnapshot
		wantResult bool
	}{
		{
			name:   "Allowed Exact Route Matching",
			roleID: 1,
			method: "GET",
			path:   "/api/v1/me",
			snapshot: rolePermissionsSnapshot{
				allRules: []domain.Permission{
					{ID: 1, Method: "GET", Path: "/api/v1/me"},
				},
				roleRules: map[int64][]domain.Permission{
					1: {
						{ID: 1, Method: "GET", Path: "/api/v1/me"},
					},
				},
			},
			wantResult: true,
		},
		{
			name:   "Denied Non-Registered Route",
			roleID: 1,
			method: "POST",
			path:   "/api/v1/unknown",
			snapshot: rolePermissionsSnapshot{
				allRules: []domain.Permission{
					{ID: 1, Method: "GET", Path: "/api/v1/me"},
				},
				roleRules: map[int64][]domain.Permission{
					1: {
						{ID: 1, Method: "GET", Path: "/api/v1/me"},
					},
				},
			},
			wantResult: false,
		},
		{
			name:   "Denied Unassigned Route",
			roleID: 2,
			method: "GET",
			path:   "/api/v1/me",
			snapshot: rolePermissionsSnapshot{
				allRules: []domain.Permission{
					{ID: 1, Method: "GET", Path: "/api/v1/me"},
				},
				roleRules: map[int64][]domain.Permission{
					1: {
						{ID: 1, Method: "GET", Path: "/api/v1/me"},
					},
				},
			},
			wantResult: false,
		},
		{
			name:   "Allowed Route with Param Wildcard Match",
			roleID: 1,
			method: "PUT",
			path:   "/api/v1/roles/123",
			snapshot: rolePermissionsSnapshot{
				allRules: []domain.Permission{
					{ID: 1, Method: "PUT", Path: "/api/v1/roles/:id"},
				},
				roleRules: map[int64][]domain.Permission{
					1: {
						{ID: 1, Method: "PUT", Path: "/api/v1/roles/:id"},
					},
				},
			},
			wantResult: true,
		},
		{
			name:   "Denied Param Wildcard Mismatch Path",
			roleID: 1,
			method: "PUT",
			path:   "/api/v1/roles/123/extra",
			snapshot: rolePermissionsSnapshot{
				allRules: []domain.Permission{
					{ID: 1, Method: "PUT", Path: "/api/v1/roles/:id"},
				},
				roleRules: map[int64][]domain.Permission{
					1: {
						{ID: 1, Method: "PUT", Path: "/api/v1/roles/:id"},
					},
				},
			},
			wantResult: false,
		},
		{
			name:   "Allowed Case Insensitive Method Matching",
			roleID: 1,
			method: "get",
			path:   "/api/v1/me",
			snapshot: rolePermissionsSnapshot{
				allRules: []domain.Permission{
					{ID: 1, Method: "GET", Path: "/api/v1/me"},
				},
				roleRules: map[int64][]domain.Permission{
					1: {
						{ID: 1, Method: "GET", Path: "/api/v1/me"},
					},
				},
			},
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &RBACService{
				snapshot: tt.snapshot,
			}
			got := svc.Allowed(tt.roleID, tt.method, tt.path)
			if got != tt.wantResult {
				t.Errorf("RBACService.Allowed() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
