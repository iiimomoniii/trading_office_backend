package handler

import (
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// Live godoc
// @Summary Health check
// @Tags dashboard
// @Produce json
// @Success 200
// @Router /live [get]
func (h *DashboardHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"version": "0.1.0",
	})
}
