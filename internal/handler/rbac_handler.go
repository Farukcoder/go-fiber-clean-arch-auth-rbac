package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/dto"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/service"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

type RBACHandler struct {
	service *service.RBACService
}

func NewRBACHandler(service *service.RBACService) *RBACHandler {
	return &RBACHandler{service: service}
}

func (h *RBACHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.service.ListRoles(context.Background())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse(http.StatusInternalServerError, "failed to retrieve roles", nil))
	}
	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Roles retrieved successfully", roles))
}

func (h *RBACHandler) CreateRole(c *fiber.Ctx) error {
	var input dto.CreateRoleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil))
	}

	role, err := h.service.CreateRole(context.Background(), domain.Role{Name: input.Name, Description: input.Description})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusCreated).JSON(dto.SuccessResponse(http.StatusCreated, "Role created successfully", role))
}

func (h *RBACHandler) UpdateRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid role id", nil))
	}

	var input dto.UpdateRoleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil))
	}

	role, err := h.service.UpdateRole(context.Background(), roleID, domain.Role{Name: input.Name, Description: input.Description})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Role updated successfully", role))
}

func (h *RBACHandler) DeleteRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid role id", nil))
	}

	if err := h.service.DeleteRole(context.Background(), roleID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Role deleted successfully", nil))
}

func (h *RBACHandler) ListPermissions(c *fiber.Ctx) error {
	permissions, err := h.service.ListPermissions(context.Background())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse(http.StatusInternalServerError, "failed to retrieve permissions", nil))
	}
	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permissions retrieved successfully", permissions))
}

func (h *RBACHandler) CreatePermission(c *fiber.Ctx) error {
	var input dto.CreatePermissionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil))
	}

	permission, err := h.service.CreatePermission(context.Background(), domain.Permission{Name: input.Name, Module: input.Module, Method: input.Method, Path: input.Path, Description: input.Description})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusCreated).JSON(dto.SuccessResponse(http.StatusCreated, "Permission created successfully", permission))
}

func (h *RBACHandler) UpdatePermission(c *fiber.Ctx) error {
	permissionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid permission id", nil))
	}

	var input dto.UpdatePermissionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil))
	}

	permission, err := h.service.UpdatePermission(context.Background(), permissionID, domain.Permission{Name: input.Name, Module: input.Module, Method: input.Method, Path: input.Path, Description: input.Description})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permission updated successfully", permission))
}

func (h *RBACHandler) DeletePermission(c *fiber.Ctx) error {
	permissionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid permission id", nil))
	}

	if err := h.service.DeletePermission(context.Background(), permissionID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permission deleted successfully", nil))
}

func (h *RBACHandler) AssignPermissionToRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid role id", nil))
	}

	var input dto.AssignPermissionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil))
	}

	if err := h.service.AssignPermissionToRole(context.Background(), roleID, input.PermissionID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permission assigned to role successfully", nil))
}

func (h *RBACHandler) RevokePermissionFromRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid role id", nil))
	}

	permissionID, err := strconv.ParseInt(c.Params("permission_id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid permission id", nil))
	}

	if err := h.service.RevokePermissionFromRole(context.Background(), roleID, permissionID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permission revoked from role successfully", nil))
}

func (h *RBACHandler) AssignRoleToUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid user id", nil))
	}

	var input dto.AssignRoleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid request payload", nil))
	}

	if err := h.service.AssignRoleToUser(context.Background(), userID, input.RoleID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Role assigned to user successfully", nil))
}

func (h *RBACHandler) GetUserPermissions(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(dto.ErrorResponse(http.StatusUnauthorized, "unauthorized", nil))
	}

	claims := token.Claims.(jwt.MapClaims)
	roleIDFloat, ok := claims["role_id"].(float64)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid role id in token", nil))
	}

	permissions, err := h.service.GetRolePermissions(context.Background(), int64(roleIDFloat))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse(http.StatusInternalServerError, "failed to retrieve permissions", nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permissions retrieved successfully", permissions))
}

func (h *RBACHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse(http.StatusBadRequest, "invalid role id", nil))
	}

	permissions, err := h.service.GetRolePermissions(context.Background(), roleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse(http.StatusInternalServerError, "failed to retrieve permissions", nil))
	}

	return c.Status(http.StatusOK).JSON(dto.SuccessResponse(http.StatusOK, "Permissions retrieved successfully", permissions))
}
