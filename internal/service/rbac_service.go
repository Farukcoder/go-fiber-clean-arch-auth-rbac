package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/repository"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

type rolePermissionsSnapshot struct {
	roles       map[int64]domain.Role
	permissions map[int64]domain.Permission
	roleRules   map[int64][]domain.Permission
	allRules    []domain.Permission
}

type RBACService struct {
	repo     *repository.RBACRepository
	mu       sync.RWMutex
	snapshot rolePermissionsSnapshot
}

func NewRBACService(repo *repository.RBACRepository) *RBACService {
	return &RBACService{repo: repo}
}

func (s *RBACService) Reload(ctx context.Context) error {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return err
	}

	permissions, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return err
	}

	details, err := s.repo.LoadRolePermissionDetails(ctx)
	if err != nil {
		return err
	}

	snapshot := rolePermissionsSnapshot{
		roles:       make(map[int64]domain.Role, len(roles)),
		permissions: make(map[int64]domain.Permission, len(permissions)),
		roleRules:   make(map[int64][]domain.Permission),
		allRules:    make([]domain.Permission, 0, len(permissions)),
	}

	for _, role := range roles {
		snapshot.roles[role.ID] = role
	}

	for _, permission := range permissions {
		snapshot.permissions[permission.ID] = permission
		snapshot.allRules = append(snapshot.allRules, permission)
	}

	for _, detail := range details {
		permission := domain.Permission{
			ID:          detail.PermissionID,
			Name:        detail.PermissionName,
			Method:      strings.ToUpper(detail.PermissionMethod),
			Path:        detail.PermissionPath,
			Description: detail.PermissionDesc,
		}
		snapshot.roleRules[detail.RoleID] = append(snapshot.roleRules[detail.RoleID], permission)
	}

	s.mu.Lock()
	s.snapshot = snapshot
	s.mu.Unlock()

	return nil
}

func (s *RBACService) StartAutoReload(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		return
	}

	ticker := time.NewTicker(interval)
	go func() { // #nosec G118 — intentional background reload, not request-scoped
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = s.Reload(context.Background())
			}
		}
	}()
}

func (s *RBACService) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := strings.ToUpper(c.Method())
		path := c.Path()

		// Step 1: extract jwt_claims set by SuccessHandler in routes.go
		rawClaims := c.Locals("jwt_claims")
		claims, ok := rawClaims.(jwt.MapClaims)
		if !ok {
			slog.Warn("RBAC: jwt_claims missing or wrong type",
				"path", path,
				"method", method,
				"raw_type", fmt.Sprintf("%T", rawClaims),
			)
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}

		// Step 2: extract role_id from claims
		roleID, err := claimAsInt64(claims["role_id"])
		if err != nil || roleID == 0 {
			slog.Warn("RBAC: role_id missing or zero in claims",
				"path", path,
				"method", method,
				"role_id_raw", claims["role_id"],
				"error", err,
			)
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}

		// Step 3: check permission
		if s.Allowed(roleID, method, path) {
			return c.Next()
		}

		s.mu.RLock()
		snap := s.snapshot
		s.mu.RUnlock()
		slog.Warn("RBAC: access denied",
			"path", path,
			"method", method,
			"role_id", roleID,
			"snapshot_roles", len(snap.roles),
			"snapshot_permissions", len(snap.permissions),
			"role_permissions", len(snap.roleRules[roleID]),
			"all_rules", len(snap.allRules),
		)
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}
}

func (s *RBACService) Allowed(roleID int64, method string, actualPath string) bool {
	s.mu.RLock()
	snapshot := s.snapshot
	s.mu.RUnlock()

	method = strings.ToUpper(method)

	if !routeRegistered(snapshot.allRules, method, actualPath) {
		slog.Debug("RBAC: route not registered in allRules",
			"method", method,
			"path", actualPath,
			"total_rules", len(snapshot.allRules),
		)
		return false
	}

	permissions := snapshot.roleRules[roleID]
	slog.Debug("RBAC: checking role permissions",
		"role_id", roleID,
		"method", method,
		"path", actualPath,
		"role_permission_count", len(permissions),
	)
	for _, permission := range permissions {
		if strings.EqualFold(permission.Method, method) && matchPath(permission.Path, actualPath) {
			return true
		}
	}

	return false
}

func (s *RBACService) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *RBACService) CreateRole(ctx context.Context, input domain.Role) (*domain.Role, error) {
	role, err := s.repo.CreateRole(ctx, input)
	if err != nil {
		return nil, err
	}
	return role, s.Reload(ctx)
}

func (s *RBACService) UpdateRole(ctx context.Context, id int64, input domain.Role) (*domain.Role, error) {
	role, err := s.repo.UpdateRole(ctx, id, input)
	if err != nil {
		return nil, err
	}
	return role, s.Reload(ctx)
}

func (s *RBACService) DeleteRole(ctx context.Context, id int64) error {
	if err := s.repo.DeleteRole(ctx, id); err != nil {
		return err
	}
	return s.Reload(ctx)
}

func (s *RBACService) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *RBACService) CreatePermission(ctx context.Context, input domain.Permission) (*domain.Permission, error) {
	permission, err := s.repo.CreatePermission(ctx, input)
	if err != nil {
		return nil, err
	}
	return permission, s.Reload(ctx)
}

func (s *RBACService) UpdatePermission(ctx context.Context, id int64, input domain.Permission) (*domain.Permission, error) {
	permission, err := s.repo.UpdatePermission(ctx, id, input)
	if err != nil {
		return nil, err
	}
	return permission, s.Reload(ctx)
}

func (s *RBACService) DeletePermission(ctx context.Context, id int64) error {
	if err := s.repo.DeletePermission(ctx, id); err != nil {
		return err
	}
	return s.Reload(ctx)
}

func (s *RBACService) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	if err := s.repo.AssignPermissionToRole(ctx, roleID, permissionID); err != nil {
		return err
	}
	return s.Reload(ctx)
}

func (s *RBACService) RevokePermissionFromRole(ctx context.Context, roleID int64, permissionID int64) error {
	if err := s.repo.RevokePermissionFromRole(ctx, roleID, permissionID); err != nil {
		return err
	}
	return s.Reload(ctx)
}

func (s *RBACService) AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error {
	if err := s.repo.AssignRoleToUser(ctx, userID, roleID); err != nil {
		return err
	}
	return s.Reload(ctx)
}

func (s *RBACService) GetRolePermissions(ctx context.Context, roleID int64) ([]domain.Permission, error) {
	// Fetch permissions directly from database to ensure up-to-date data
	permissions, err := s.repo.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		slog.Error("RBAC: Failed to fetch permissions from database", "role_id", roleID, "error", err)
		return []domain.Permission{}, err
	}

	slog.Info("RBAC: GetRolePermissions returned from database",
		"role_id", roleID,
		"count", len(permissions),
	)

	// If no permissions found in database, return empty array
	if len(permissions) == 0 {
		slog.Warn("RBAC: No permissions found for role in database", "role_id", roleID)
		return []domain.Permission{}, nil
	}

	return permissions, nil
}

func (s *RBACService) getDefaultPermissionsForRole(roleName string) []domain.Permission {
	// Define default permissions for each role type
	defaultPermissions := map[string][]string{
		"superadmin": {
			"view_dashboard", "view_products", "view_categories", "view_subcategories",
			"view_stock", "manage_roles", "manage_permissions", "manage_user_roles", "manage_settings",
		},
		"admin": {
			"view_dashboard", "view_products", "view_categories", "view_subcategories",
			"view_stock", "manage_roles", "manage_permissions", "manage_user_roles", "manage_settings",
		},
		"manager": {
			"view_dashboard", "view_products", "view_categories", "view_subcategories",
			"view_stock", "manage_settings",
		},
		"staff": {
			"view_dashboard", "view_products", "view_stock",
		},
		"customer": {
			"view_dashboard",
		},
	}

	permissionNames, ok := defaultPermissions[roleName]
	if !ok {
		// If role not found, return minimal permissions
		return []domain.Permission{}
	}

	s.mu.RLock()
	snapshot := s.snapshot
	s.mu.RUnlock()

	// Convert permission names to Permission objects from the snapshot
	var result []domain.Permission
	permissionNameSet := make(map[string]bool)
	for _, name := range permissionNames {
		permissionNameSet[name] = true
	}

	for _, permission := range snapshot.allRules {
		if permissionNameSet[permission.Name] {
			result = append(result, permission)
		}
	}

	return result
}

func routeRegistered(permissions []domain.Permission, method string, actualPath string) bool {
	for _, permission := range permissions {
		if strings.EqualFold(permission.Method, method) && matchPath(permission.Path, actualPath) {
			return true
		}
	}
	return false
}

func matchPath(pattern string, actual string) bool {
	pattern = normalizePath(pattern)
	actual = normalizePath(actual)

	if pattern == actual {
		return true
	}

	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\:`, `:`)
	regexPattern = regexp.MustCompile(`:(\w+)`).ReplaceAllString(regexPattern, `[^/]+`)
	regexPattern = strings.ReplaceAll(regexPattern, `\*`, `.*`)

	re := regexp.MustCompile("^" + regexPattern + "$")
	return re.MatchString(actual)
}

func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	return path
}

func claimAsInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		var parsed int64
		_, err := fmt.Sscan(v, &parsed)
		return parsed, err
	default:
		return 0, errors.New("invalid claim type")
	}
}

func sortedPermissionNames(permissions []domain.Permission) []string {
	names := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		names = append(names, permission.Name)
	}
	sort.Strings(names)
	return names
}
