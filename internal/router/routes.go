package router

import (
	"go-fiber-clean-arch-auth-rbac/internal/config"
	"go-fiber-clean-arch-auth-rbac/internal/handler"
	"go-fiber-clean-arch-auth-rbac/internal/middleware"
	"go-fiber-clean-arch-auth-rbac/internal/service"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func Setup(app *fiber.App, authHandler *handler.AuthHandler, logsHandler *handler.LogsHandler, rbacHandler *handler.RBACHandler, rbacService *service.RBACService, cfg *config.Config) {
	api := app.Group("/api/v1")
	auth := api.Group("/auth", middleware.AuthRateLimiter())
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	protected := api.Group("", jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(cfg.JwtSecret)},
		ContextKey:   "user",
		ErrorHandler: middleware.ErrorHandler,
		SuccessHandler: func(c *fiber.Ctx) error {
			token, ok := c.Locals("user").(*jwt.Token)
			if ok && token != nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					c.Locals("jwt_claims", claims)
				}
			}
			return c.Next()
		},
	}))
	protected.Use(rbacService.Middleware())
	protected.Get("/me", authHandler.Me)
	protected.Get("/logs", logsHandler.GetAll)

	protected.Get("/roles", rbacHandler.ListRoles)
	protected.Post("/roles", rbacHandler.CreateRole)
	protected.Put("/roles/:id", rbacHandler.UpdateRole)
	protected.Delete("/roles/:id", rbacHandler.DeleteRole)
	protected.Get("/permissions", rbacHandler.ListPermissions)
	protected.Post("/permissions", rbacHandler.CreatePermission)
	protected.Put("/permissions/:id", rbacHandler.UpdatePermission)
	protected.Delete("/permissions/:id", rbacHandler.DeletePermission)
	protected.Post("/roles/:id/permissions", rbacHandler.AssignPermissionToRole)
	protected.Delete("/roles/:id/permissions/:permission_id", rbacHandler.RevokePermissionFromRole)
	protected.Patch("/users/:id/role", rbacHandler.AssignRoleToUser)
}
