package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

type CreateUserInput struct {
	Name         string
	Email        string
	Phone        string
	PasswordHash string
	RoleID       int64
}

type userRecord struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	Name         string `gorm:"column:name"`
	Email        string `gorm:"column:email"`
	Phone        string `gorm:"column:phone"`
	PasswordHash string `gorm:"column:password_hash"`
}

func buildUserRecord(input CreateUserInput) userRecord {
	return userRecord{
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: input.PasswordHash,
	}
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.email, users.phone, COALESCE(ur.role_id, 0) AS role_id, COALESCE(r.name, '') AS role_name, users.password_hash, users.created_at").
		Joins("LEFT JOIN user_roles ur ON ur.user_id = users.id").
		Joins("LEFT JOIN roles r ON r.id = ur.role_id").
		Where("users.email = ?", strings.TrimSpace(email)).
		Take(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.email, users.phone, COALESCE(ur.role_id, 0) AS role_id, COALESCE(r.name, '') AS role_name, users.password_hash, users.created_at").
		Joins("LEFT JOIN user_roles ur ON ur.user_id = users.id").
		Joins("LEFT JOIN roles r ON r.id = ur.role_id").
		Where("users.id = ?", id).
		Take(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.email, users.phone, COALESCE(ur.role_id, 0) AS role_id, COALESCE(r.name, '') AS role_name, users.password_hash, users.created_at").
		Joins("LEFT JOIN user_roles ur ON ur.user_id = users.id").
		Joins("LEFT JOIN roles r ON r.id = ur.role_id").
		Where("users.phone = ?", strings.TrimSpace(phone)).
		Take(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	record := buildUserRecord(input)

	if err := r.db.WithContext(ctx).Table("users").Create(&record).Error; err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           record.ID,
		Name:         record.Name,
		Email:        record.Email,
		Phone:        record.Phone,
		PasswordHash: record.PasswordHash,
	}

	if input.RoleID > 0 {
		now := time.Now().UTC()
		if err := r.db.WithContext(ctx).Table("user_roles").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{"role_id": input.RoleID, "updated_at": now}),
		}).Create(&domain.UserRole{UserID: user.ID, RoleID: input.RoleID, CreatedAt: now, UpdatedAt: now}).Error; err != nil {
			return nil, err
		}
		user.RoleID = input.RoleID

		var role domain.Role
		if err := r.db.WithContext(ctx).Table("roles").Select("name").Where("id = ?", input.RoleID).Take(&role).Error; err == nil {
			user.RoleName = role.Name
		}
	}

	return user, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID int64, roleID int64) error {
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

func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.email, users.phone, COALESCE(ur.role_id, 0) AS role_id, COALESCE(r.name, '') AS role_name, users.created_at").
		Joins("LEFT JOIN user_roles ur ON ur.user_id = users.id").
		Joins("LEFT JOIN roles r ON r.id = ur.role_id").
		Order("users.created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
