package middleware

import (
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/service"
	"github.com/gofiber/fiber/v2"
)

func RBAC(svc *service.RBACService) fiber.Handler {
	return svc.Middleware()
}
