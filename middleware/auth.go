package middleware

import (
	"github.com/gofiber/fiber/v2"

	"trading-office/trading_office_backend/utils"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized,
				"GLB004", "UNAUTHORIZED",
				"Unauthorized",
				"ไม่มีสิทธิ์เข้าถึง",
			)
		}
		return c.Next()
	}
}
