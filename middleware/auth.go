package middleware

import (
	"errors"
	"strings"

	"beli_mang/crypto"
	"beli_mang/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return responses.NewUnauthorizedError("token not found")
	}

	splitted := strings.Split(auth, " ")
	if len(splitted) != 2 || splitted[0] != "Bearer" {
		return responses.NewUnauthorizedError("invalid token")
	}

	payload, err := crypto.VerifyToken(splitted[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return responses.NewUnauthorizedError("token expired")
		}
		return responses.NewUnauthorizedError(err.Error())
	}

	if payload.Role != "admin" {
		return responses.NewUnauthorizedError("user is not an admin")
	}

	ctx.Locals("user_id", payload.Id)
	ctx.Locals("username", payload.Username)
	ctx.Locals("role", payload.Role)

	return ctx.Next()
}

func UserAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return responses.NewUnauthorizedError("token not found")
	}

	splitted := strings.Split(auth, " ")
	if len(splitted) != 2 || splitted[0] != "Bearer" {
		return responses.NewUnauthorizedError("invalid token")
	}

	payload, err := crypto.VerifyToken(splitted[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return responses.NewUnauthorizedError("token expired")
		}
		return responses.NewUnauthorizedError(err.Error())
	}

	if payload.Role != "user" {
		return responses.NewUnauthorizedError("user is not a user")
	}

	ctx.Locals("user_id", payload.Id)
	ctx.Locals("username", payload.Username)
	ctx.Locals("role", payload.Role)

	return ctx.Next()
}

func Auth(ctx *fiber.Ctx) error {
	auth := ctx.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return responses.NewUnauthorizedError("token not found")
	}

	splitted := strings.Split(auth, " ")
	if len(splitted) != 2 || splitted[0] != "Bearer" {
		return responses.NewUnauthorizedError("invalid token")
	}

	payload, err := crypto.VerifyToken(splitted[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return responses.NewUnauthorizedError("token expired")
		}
		return responses.NewUnauthorizedError(err.Error())
	}

	ctx.Locals("user_id", payload.Id)
	ctx.Locals("username", payload.Username)
	ctx.Locals("role", payload.Role)

	return ctx.Next()
}
