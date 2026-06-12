package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"trading-office/trading_office_backend/model"
)

func ErrorResponse(c *fiber.Ctx, status int, errorCode, errName, messageEN, messageTH string) error {
	return c.Status(status).JSON(model.BaseResponse{
		Status:  status,
		Message: messageEN,
		Data: model.ErrorDetail{
			Timestamp: time.Now(),
			Status:    status,
			ErrorCode: errorCode,
			Error:     errName,
			MessageEN: messageEN,
			MessageTH: messageTH,
			Path:      c.Path(),
		},
	})
}

func SuccessResponse(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(model.BaseResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func TooManyRequests(c *fiber.Ctx) error {
	return ErrorResponse(c, fiber.StatusTooManyRequests,
		"GLB005", "TOO_MANY_REQUESTS",
		"Too many requests, please try again later",
		"คำขอมากเกินไป กรุณาลองใหม่ภายหลัง",
	)
}
