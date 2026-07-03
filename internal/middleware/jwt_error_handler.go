package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTErrorHandler(c *fiber.Ctx) error {
	return c.Next()
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	errorType := "UNKNOWN_ERROR"
	isJWTError := false

	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "missing") || strings.Contains(errMsg, "malformed") {
		errorType = "JWT_MISSING_OR_MALFORMED"
		isJWTError = true
	} else if strings.Contains(errMsg, "expired") || strings.Contains(errMsg, "invalid") {
		errorType = "JWT_INVALID_OR_EXPIRED"
		isJWTError = true
	} else {
		errorType = "INTERNAL_SERVER_ERROR"
	}

	slog.Error("Request error",
		"error_type", errorType,
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	if isJWTError {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or expired JWT",
			"data":    nil,
		})
	}

	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"status":  "error",
		"message": "Internal server error",
		"data":    nil,
	})
}
