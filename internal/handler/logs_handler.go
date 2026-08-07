package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/dto"
	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type LogsHandler struct {
	requestLogRepo *repository.RequestLogRepository
}

func NewLogsHandler(requestLogRepo *repository.RequestLogRepository) *LogsHandler {
	return &LogsHandler{requestLogRepo: requestLogRepo}
}

func (h *LogsHandler) GetAll(c *fiber.Ctx) error {
	methodFilter := c.Query("method", "")
	statusCodeFilter := 0
	if statusCodeStr := c.Query("status_code", ""); statusCodeStr != "" {
		if v, err := strconv.Atoi(statusCodeStr); err == nil {
			statusCodeFilter = v
		}
	}

	page := 1
	if pageStr := c.Query("page", "1"); pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	perPage := 15
	if perPageStr := c.Query("per_page", "15"); perPageStr != "" {
		if v, err := strconv.Atoi(perPageStr); err == nil && v > 0 {
			perPage = v
		}
	}

	logs, total, err := h.requestLogRepo.GetAll(context.Background(), methodFilter, statusCodeFilter, page, perPage)
	if err != nil {
		response := dto.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve logs", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	totalPages := (total + perPage - 1) / perPage
	hasNext := page < totalPages
	hasPrev := page > 1

	responseData := map[string]interface{}{
		"data": logs,
		"pagination": dto.PaginationResponse{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}

	response := dto.SuccessResponse(http.StatusOK, "Logs retrieved successfully", responseData)
	return c.Status(http.StatusOK).JSON(response)
}
