package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SecurityHeaders returns a middleware that injects security headers.
func SecurityHeaders(appEnv string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Enable HSTS in production/staging environments
		env := strings.ToLower(appEnv)
		if env != "local" && env != "dev" && env != "development" && env != "" {
			c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		return c.Next()
	}
}

// CORS returns a middleware configured with explicit allowed origins.
func CORS(allowedOrigins string) fiber.Handler {
	if allowedOrigins == "" || allowedOrigins == "*" {
		// Fallback defaults for safety, but typically overridden by config
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:5174"
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	})
}
