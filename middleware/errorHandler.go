package middleware

import (
	"beli_mang/responses"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if customError, ok := err.(responses.CustomError); ok {
		code = customError.Status()
		return ctx.Status(code).JSON(map[string]interface{}{
			"message": customError.Error(),
		})
	}

	return ctx.Status(code).JSON(map[string]interface{}{
		"message": err.Error(),
	})
}
