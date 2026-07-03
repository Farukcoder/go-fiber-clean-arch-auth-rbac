package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"go-fiber-clean-arch-auth-rbac/internal/dto"
	"go-fiber-clean-arch-auth-rbac/internal/service"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
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

	response := dto.SuccessResponse(http.StatusOK, "User retrieved successfully", user)
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
