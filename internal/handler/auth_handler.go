package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/dto"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/service"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service     *service.AuthService
	rbacService *service.RBACService
}

func NewAuthHandler(service *service.AuthService, rbacService *service.RBACService) *AuthHandler {
	return &AuthHandler{service: service, rbacService: rbacService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input dto.LoginInput
	if err := c.BodyParser(&input); err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	resp, err := h.service.Login(context.Background(), input)
	if err != nil {
		response := dto.ErrorResponse(http.StatusUnauthorized, err.Error(), nil)
		return c.Status(http.StatusUnauthorized).JSON(response)
	}

	response := dto.SuccessResponse(http.StatusOK, "Login successful", resp)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Validate password confirmation
	if input.Password != input.ConfirmPassword {
		response := dto.ErrorResponse(http.StatusBadRequest, "password and confirm password do not match", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	resp, err := h.service.Register(context.Background(), input)
	if err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	response := dto.SuccessResponse(http.StatusCreated, "Registration successful", resp)
	return c.Status(http.StatusCreated).JSON(response)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var input dto.RefreshRequest
	if err := c.BodyParser(&input); err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	if input.RefreshToken == "" {
		response := dto.ErrorResponse(http.StatusBadRequest, "refresh_token is required", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	resp, err := h.service.Refresh(context.Background(), input)
	if err != nil {
		response := dto.ErrorResponse(http.StatusUnauthorized, err.Error(), nil)
		return c.Status(http.StatusUnauthorized).JSON(response)
	}

	response := dto.SuccessResponse(http.StatusOK, "Token refreshed successfully", resp)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var input dto.RefreshRequest
	if err := c.BodyParser(&input); err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	if input.RefreshToken == "" {
		response := dto.ErrorResponse(http.StatusBadRequest, "refresh_token is required", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	if err := h.service.Logout(context.Background(), input); err != nil {
		response := dto.ErrorResponse(http.StatusUnauthorized, err.Error(), nil)
		return c.Status(http.StatusUnauthorized).JSON(response)
	}

	response := dto.SuccessResponse(http.StatusOK, "Logged out successfully", nil)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// "jwt_claims" is set by the JWT SuccessHandler in routes.go as jwtv4.MapClaims
	claims, ok := c.Locals("jwt_claims").(jwt.MapClaims)
	if !ok {
		response := dto.ErrorResponse(http.StatusUnauthorized, "Unauthorized", nil)
		return c.Status(http.StatusUnauthorized).JSON(response)
	}

	userID, err := claimUserID(claims)
	if err != nil {
		response := dto.ErrorResponse(http.StatusUnauthorized, "Invalid token claims", nil)
		return c.Status(http.StatusUnauthorized).JSON(response)
	}

	user, err := h.service.GetUserByID(context.Background(), int(userID))
	if err != nil {
		response := dto.ErrorResponse(http.StatusNotFound, "User not found", nil)
		return c.Status(http.StatusNotFound).JSON(response)
	}

	// Get user's role permissions
	roleID, err := claimAsInt64(claims["role_id"])
	if err != nil {
		roleID = user.RoleID
	}

	slog.Info("Fetching permissions for user", "user_id", userID, "role_id", roleID)

	permissions, err := h.rbacService.GetRolePermissions(context.Background(), roleID)
	if err != nil {
		slog.Error("Error fetching permissions", "role_id", roleID, "error", err)
		// If permissions can't be loaded, return empty array
		permissions = []domain.Permission{}
	}

	slog.Info("Permissions fetched", "role_id", roleID, "count", len(permissions))

	// Combine user data with permissions
	userData := map[string]interface{}{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"phone":       user.Phone,
		"role_id":     user.RoleID,
		"role_name":   user.RoleName,
		"permissions": permissions,
	}

	response := dto.SuccessResponse(http.StatusOK, "User retrieved successfully", userData)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers(context.Background())
	if err != nil {
		response := dto.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve users", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := dto.SuccessResponse(http.StatusOK, "Users retrieved successfully", users)
	return c.Status(http.StatusOK).JSON(response)
}

func claimUserID(claims jwt.MapClaims) (int64, error) {
	if value, ok := claims["user_id"]; ok {
		return claimAsInt64(value)
	}
	if value, ok := claims["sub"]; ok {
		return claimAsInt64(value)
	}
	return 0, errors.New("invalid token claims")
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
		parsed, err := strconv.ParseInt(v, 10, 64)
		return parsed, err
	default:
		return 0, errors.New("invalid claim type")
	}
}
