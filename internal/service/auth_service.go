package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"go-fiber-clean-arch-auth-rbac/internal/domain"
	"go-fiber-clean-arch-auth-rbac/internal/dto"
	"go-fiber-clean-arch-auth-rbac/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL   = 15 * time.Minute
	refreshTokenTTL  = 2 * 24 * time.Hour
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	rbacRepo         *repository.RBACRepository
	refreshRepo      *repository.RefreshTokenRepository
	jwtSecret        string
	jwtRefreshSecret string
}

func NewAuthService(userRepo *repository.UserRepository, rbacRepo *repository.RBACRepository, refreshRepo *repository.RefreshTokenRepository, jwtSecret string, jwtRefreshSecret string) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		rbacRepo:         rbacRepo,
		refreshRepo:      refreshRepo,
		jwtSecret:        jwtSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
}

func (s *AuthService) Login(ctx context.Context, input dto.LoginInput) (*dto.TokenResponse, error) {
	user, err := s.findUserByIdentifier(ctx, input.EmailOrPhone)
	if err != nil {
		slog.Error("User not found", "email_or_phone", input.EmailOrPhone, "error", err)
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		slog.Error("Password comparison failed", "email_or_phone", input.EmailOrPhone, "error", err)
		return nil, errors.New("invalid credentials")
	}

	resp, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AuthService) Refresh(ctx context.Context, input dto.RefreshRequest) (*dto.TokenResponse, error) {
	claims, err := s.parseToken(input.RefreshToken, s.jwtRefreshSecret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if claims.TokenType != tokenTypeRefresh {
		return nil, errors.New("invalid refresh token")
	}

	if claims.JTI == "" {
		return nil, errors.New("invalid refresh token")
	}

	hashedJTI := s.hashJTI(claims.JTI)
	storedToken, err := s.refreshRepo.FindByJTI(ctx, hashedJTI)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if storedToken.RevokedAt.Valid {
		// Compromise detection: revoke all active tokens for this user!
		slog.Warn("Refresh token reuse detected! Revoking all sessions for user", "user_id", storedToken.UserID)
		_ = s.refreshRepo.RevokeUserTokens(ctx, storedToken.UserID)
		return nil, errors.New("invalid or expired refresh token")
	}

	if time.Now().UTC().After(storedToken.ExpiresAt.UTC()) {
		return nil, errors.New("invalid or expired refresh token")
	}

	if claims.UserID != storedToken.UserID {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, int(claims.UserID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.refreshRepo.RevokeByJTI(ctx, hashedJTI); err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, input dto.RefreshRequest) error {
	claims, err := s.parseToken(input.RefreshToken, s.jwtRefreshSecret)
	if err != nil {
		return errors.New("invalid or expired refresh token")
	}

	if claims.JTI == "" {
		return errors.New("invalid refresh token")
	}

	return s.refreshRepo.RevokeByJTI(ctx, s.hashJTI(claims.JTI))
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterRequest) (*dto.UserResponse, error) {
	slog.Info("Registration attempt", "email", input.Email, "phone", input.Phone)

	customerRole, err := s.rbacRepo.FindRoleByName(ctx, "customer")
	if err != nil {
		slog.Warn("Customer role not found during registration, falling back to zero role", "error", err)
		customerRole = &domain.Role{}
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		slog.Error("User already exists", "email", input.Email)
		return nil, errors.New("user with this email already exists")
	}

	if input.Phone != "" {
		existingUser, err := s.userRepo.FindByPhone(ctx, input.Phone)
		if err == nil && existingUser != nil {
			slog.Error("User already exists", "phone", input.Phone)
			return nil, errors.New("user with this phone already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Password hashing failed", "error", err)
		return nil, err
	}

	roleID := int64(0)
	if customerRole != nil && customerRole.ID > 0 {
		roleID = customerRole.ID
	}

	user, err := s.userRepo.Create(ctx, repository.CreateUserInput{
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		RoleID:       roleID,
	})
	if err != nil {
		slog.Error("User creation failed", "error", err)
		return nil, err
	}

	if roleID > 0 {
		user.RoleID = roleID
		if role, err := s.rbacRepo.FindRoleByID(ctx, roleID); err == nil && role != nil && role.Name != "" {
			user.RoleName = role.Name
		} else if customerRole != nil && customerRole.Name != "" {
			user.RoleName = customerRole.Name
		}
	}

	slog.Info("Registration successful", "id", user.ID, "email", user.Email)

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		RoleID:   user.RoleID,
		RoleName: user.RoleName,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		RoleID:   user.RoleID,
		RoleName: user.RoleName,
	}, nil
}

func (s *AuthService) findUserByIdentifier(ctx context.Context, identifier string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, identifier)
	if err == nil {
		return user, nil
	}

	user, err = s.userRepo.FindByPhone(ctx, identifier)
	if err == nil {
		return user, nil
	}

	return nil, err
}

func (s *AuthService) issueTokenPair(ctx context.Context, user *domain.User) (*dto.TokenResponse, error) {
	now := time.Now().UTC()
	accessExpiresAt := now.Add(accessTokenTTL)
	refreshExpiresAt := now.Add(refreshTokenTTL)
	roleName := user.RoleName
	if roleName == "" && user.RoleID > 0 {
		if role, err := s.rbacRepo.FindRoleByID(ctx, user.RoleID); err == nil && role != nil && role.Name != "" {
			roleName = role.Name
		}
	}
	if roleName == "" {
		roleName = "customer"
	}

	accessToken, err := s.signToken(s.jwtSecret, dto.AuthClaims{
		UserID:    user.ID,
		RoleID:    user.RoleID,
		RoleName:  roleName,
		Email:     user.Email,
		TokenType: tokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
		},
	})
	if err != nil {
		return nil, err
	}

	refreshJTI := uuid.NewString()
	refreshToken, err := s.signToken(s.jwtRefreshSecret, dto.AuthClaims{
		UserID:    user.ID,
		RoleID:    user.RoleID,
		RoleName:  roleName,
		Email:     user.Email,
		JTI:       refreshJTI,
		TokenType: tokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
		},
	})
	if err != nil {
		return nil, err
	}

	tokenRecord := &domain.RefreshToken{
		UserID:    user.ID,
		JTI:       s.hashJTI(refreshJTI),
		ExpiresAt: refreshExpiresAt,
	}
	if err := s.refreshRepo.Create(ctx, tokenRecord); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
	}, nil
}

func (s *AuthService) signToken(secret string, claims dto.AuthClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) parseToken(tokenString string, secret string) (*dto.AuthClaims, error) {
	claims := &dto.AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	parsedClaims, ok := token.Claims.(*dto.AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return parsedClaims, nil
}

func (s *AuthService) hashJTI(jti string) string {
	hash := sha256.Sum256([]byte(jti))
	return hex.EncodeToString(hash[:])
}
