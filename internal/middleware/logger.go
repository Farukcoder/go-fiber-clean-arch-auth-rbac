package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(requestLogRepo *repository.RequestLogRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Capture request body
		var requestBody string
		if c.Method() != "GET" && !strings.Contains(c.Path(), "/auth/") {
			requestBody = string(c.Body())
		} else if strings.Contains(c.Path(), "/auth/") {
			requestBody = "[REDACTED]"
		}

		err := c.Next()

		duration := time.Since(start)

		// Format: 2026-06-30 17:05:01 GET /api/v1/me .................................. 200 ~ 1s
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		durationStr := fmt.Sprintf("~ %s", duration.String())
		dots := " .................................."

		fmt.Printf("%s %s %s%s %d %s\n",
			timestamp,
			c.Method(),
			c.Path(),
			dots,
			c.Response().StatusCode(),
			durationStr,
		)

		// Determine error type and message
		errorType := ""
		errorMessage := ""
		if err != nil {
			errorType = "REQUEST_ERROR"
			errorMessage = err.Error()
		} else if c.Response().StatusCode() >= 400 {
			errorType = "HTTP_ERROR"
			errorMessage = fmt.Sprintf("HTTP %d", c.Response().StatusCode())
		}

		// Capture all context data before goroutine
		method := c.Method()
		path := c.Path()
		statusCode := c.Response().StatusCode()
		var responseBody string
		if !strings.Contains(path, "/auth/") {
			responseBody = string(c.Response().Body())
		} else {
			responseBody = "[REDACTED]"
		}
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")
		durationMs := float64(duration.Milliseconds())

		// Save to database asynchronously
		go func() {
			log := &domain.RequestLog{
				Method:       method,
				Path:         path,
				StatusCode:   statusCode,
				DurationMs:   durationMs,
				RequestBody:  requestBody,
				ResponseBody: responseBody,
				IPAddress:    ipAddress,
				UserAgent:    userAgent,
				ErrorType:    errorType,
				ErrorMessage: errorMessage,
			}

			if err := requestLogRepo.Create(context.Background(), log); err != nil {
				slog.Error("Failed to save request log to database", "error", err)
			}
		}()

		return err
	}
}
