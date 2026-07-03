package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-fiber-clean-arch-auth-rbac/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RBACRepository struct {
	db *gorm.DB
}

type RolePermissionDetail struct {
	RoleID           int64  `gorm:"column:role_id"`
	RoleName         string `gorm:"column:role_name"`
	PermissionID     int64  `gorm:"column:permission_id"`
	PermissionName   string `gorm:"column:permission_name"`
	PermissionMethod string `gorm:"column:permission_method"`
	PermissionPath   string `gorm:"column:permission_path"`
	PermissionDesc   string `gorm:"column:permission_desc"`
}

func NewRBACRepository(db *gorm.DB) *RBACRepository {
	return &RBACRepository{db: db}
}

func (r *RBACRepository) ListRoles(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).Table("roles").Order("id ASC").Find(&roles).Error
	return roles, err
}

func (r *RBACRepository) FindRoleByID(ctx context.Context, id int64) (*domain.Role, error) {
	role := &domain.Role{}
	err := r.db.WithContext(ctx).Table("roles").Where("id = ?", id).Take(role).Error
	return role, err
}

func (r *RBACRepository) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	role := &domain.Role{}
	err := r.db.WithContext(ctx).Table("roles").Where("name = ?", strings.TrimSpace(name)).Take(role).Error
	return role, err
}

func (r *RBACRepository) CreateRole(ctx context.Context, input domain.Role) (*domain.Role, error) {
	role := &domain.Role{Name: strings.TrimSpace(input.Name), Description: input.Description}
	err := r.db.WithContext(ctx).
		Table("roles").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.Assignments(map[string]any{"description": input.Description, "updated_at": time.Now().UTC()}),
		}).
		Create(role).Error
	return role, err
}

func (r *RBACRepository) UpdateRole(ctx context.Context, id int64, input domain.Role) (*domain.Role, error) {
	result := r.db.WithContext(ctx).Table("roles").Where("id = ?", id).
		Updates(map[string]any{"name": strings.TrimSpace(input.Name), "description": input.Description, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("role not found")
	}
	return r.FindRoleByID(ctx, id)
}

func (r *RBACRepository) DeleteRole(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Table("roles").Where("id = ?", id).Delete(&domain.Role{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

func (r *RBACRepository) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Table("permissions").Order("id ASC").Find(&permissions).Error
	return permissions, err
}

func (r *RBACRepository) FindPermissionByID(ctx context.Context, id int64) (*domain.Permission, error) {
	permission := &domain.Permission{}
	err := r.db.WithContext(ctx).Table("permissions").Where("id = ?", id).Take(permission).Error
	return permission, err
}

func (r *RBACRepository) FindPermissionByName(ctx context.Context, name string) (*domain.Permission, error) {
	permission := &domain.Permission{}
	err := r.db.WithContext(ctx).Table("permissions").Where("name = ?", strings.TrimSpace(name)).Take(permission).Error
	return permission, err
}

func (r *RBACRepository) CreatePermission(ctx context.Context, input domain.Permission) (*domain.Permission, error) {
	permission := &domain.Permission{Name: strings.TrimSpace(input.Name), Method: input.Method, Path: input.Path, Description: input.Description}
	err := r.db.WithContext(ctx).
		Table("permissions").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.Assignments(map[string]any{"method": input.Method, "path": input.Path, "description": input.Description, "updated_at": time.Now().UTC()}),
		}).
		Create(permission).Error
	return permission, err
}

func (r *RBACRepository) UpdatePermission(ctx context.Context, id int64, input domain.Permission) (*domain.Permission, error) {
	result := r.db.WithContext(ctx).Table("permissions").Where("id = ?", id).
		Updates(map[string]any{"name": strings.TrimSpace(input.Name), "method": input.Method, "path": input.Path, "description": input.Description, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("permission not found")
	}
	return r.FindPermissionByID(ctx, id)
}

func (r *RBACRepository) DeletePermission(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Table("permissions").Where("id = ?", id).Delete(&domain.Permission{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

func (r *RBACRepository) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	rp := &domain.RolePermission{RoleID: roleID, PermissionID: permissionID}
	return r.db.WithContext(ctx).
		Table("role_permissions").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(rp).Error
}

func (r *RBACRepository) RevokePermissionFromRole(ctx context.Context, roleID int64, permissionID int64) error {
	result := r.db.WithContext(ctx).Table("role_permissions").Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&domain.RolePermission{})
	return result.Error
}

func (r *RBACRepository) LoadRolePermissionDetails(ctx context.Context) ([]RolePermissionDetail, error) {
	var details []RolePermissionDetail
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.id AS role_id, r.name AS role_name, p.id AS permission_id, p.name AS permission_name, p.method AS permission_method, p.path AS permission_path, p.description AS permission_desc
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id
		JOIN permissions p ON p.id = rp.permission_id
		ORDER BY r.id ASC, p.id ASC
	`).Scan(&details).Error
	return details, err
}

func (r *RBACRepository) AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error {
	var user domain.User
	if err := r.db.WithContext(ctx).Table("users").Where("id = ?", userID).Select("id").Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).Table("user_roles").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{"role_id": roleID, "updated_at": now}),
	}).Create(&domain.UserRole{UserID: userID, RoleID: roleID, CreatedAt: now, UpdatedAt: now}).Error; err != nil {
		return err
	}

	return nil
}
